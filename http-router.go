// Сервис-роутер для обмена данными с Minio Object Storage.
// TODO для прода, закоммитить отладочные методы стандартного вывода и тайминги
//
// Перед сборкой выполнить и обновить пакеты...
// $ go install golang.org/x/vuln/cmd/govulncheck@latest
// $ govulncheck ./...
// go mod tidy
//
// Проверить на известные уязвимости...
// $ go install github.com/securego/gosec/v2/cmd/gosec@latest
// $ gosec ./...
//
// [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production. Перед деплоем включить релиз.
// - using env:   export GIN_MODE=release
// - using code:  gin.SetMode(gin.ReleaseMode)
//
// TODO init flags уровня пакета для параметров подключения
// TODO обработчик сообщений для фронта. Status(500)

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Mu              sync.Mutex
	Endpoint        string `json:"s3_storage_endpoint"`
	AccessKeyID     string `json:"s3_access_key_id"`
	SecretAccessKey string `json:"s3_secret_access_key"`
	UseSSL          bool   `json:"s3_connect_use_ssl"`
	Region          string `json:"s3_storage_region"`
	ContentType     string `json:"s3_content_type"`
	UploaDir        string `json:"rout_dir_upload"`
	DownloaDir      string `json:"rout_dir_download"`
	//NameBucket      string `json:"s3_storage_name_bucket"` // uuid of user
}

func main() {

	// Switching mode. Переключение режимов отладка-релиз
	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	// router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to web-server for minio!")
	})

	// Route to down. Роутер загрузки
	router.GET("/stream/:filename", func(c *gin.Context) {

		// Get filename
		filename := c.Param("filename")
		if filename == "" {
			log.Println("Filename is empty!")
			c.Status(400)
		}

		// Get uuid
		uuid := c.GetHeader("user-uuid")
		if uuid == "" {
			log.Println("Uuid is empty!")
			c.Status(400)
		}

		var wg sync.WaitGroup // Счетчик количества горутин

		var mu sync.Mutex // Мьютекс каналов
		chs := make(chan string, 1)
		chid := make(chan string, 1)
		done := make(chan struct{})

		// Блокирует доступ к каналу, на время передачи данных в него
		mu.Lock()
		defer mu.Unlock()
		chs <- filename
		chid <- uuid

		wg.Add(1)
		start := time.Now()
		go func() { // Вызов функции создания бакета
			defer wg.Done()

			// Вызов метода создания бакета
			err := makerBucket(chid)
			if err != nil { // TODO обработчик сообщений для фронта
				c.Status(500) // Ошибка сервиса

				wg.Add(1)
				go func() { // Скачиваем из существующего бакета со значением uuid
					defer wg.Done()

					// Проверяем наличие файла локально
					file, err := os.Open("download/" + filename)
					defer file.Close()
					if err != nil { // Вызываем загрузчик файлов с бакета
						c.String(http.StatusNotFound, "File not found.")

						c.Header("Content-Type", "application/octet-stream")

						for range <-chid {
							//nil опустошаем канал со строкой сообщения api
						}
						chid <- uuid // Передаем в него uuid

						// Вызов метода загрузчика из бакета
						if fd := downloader(chs, chid); fd == "" {
							log.Println("Error download file")
							c.Status(404) // Запрошенный файл не скачен с бакета
						} else {
							log.Println("File downloaded from bucket")
							c.Status(200) // Запрошенный файл скачен с бакета
						}
					} else {
						log.Println("File is locally")
						c.Status(201) // Файл существует локально

						done <- struct{}{}
						for range <-chs {
							//nil
						}
						for range <-chid {
							//nil
						}
					}
				}()
				<-done
				go func() { // Ожидание счетчика
					wg.Wait()
					close(chs) // Закрываем канал, собираем горутины
					close(chid)
					close(done)
				}()
			}

			done <- struct{}{}
			for range <-chid {
				//nil
			}
		}()

		secs := time.Since(start).Seconds()
		log.Printf("%.2fs Time of downloader request\n", secs)

		ch1 := make(chan bool, 1)
		ch2 := make(chan bool, 1)
		ch3 := make(chan bool, 1)
		ch4 := make(chan bool, 1)

		if fm := strings.Contains(filename, ".mp4"); fm == true {
			ch1 <- true
		}
		if fm := strings.Contains(filename, ".avi"); fm == true {
			ch2 <- true
		}
		if fm := strings.Contains(filename, ".mkv"); fm == true {
			ch3 <- true
		}
		if fm := strings.Contains(filename, ".gif"); fm == true {
			ch4 <- true
		}

		// Selector of format. Селектор медиа форматов
		select {

		case <-ch1:
			wg.Add(1)
			start := time.Now()
			done := make(chan int, 1)

			go func() {
				defer wg.Done()
				file, err := os.Open("download/" + filename)
				if err != nil {
					c.String(http.StatusNotFound, "Video not found.")
					return
				}
				defer file.Close()

				c.Header("Content-Type", "video/mp4")

				// Буфер байтов переменного размера, готовый к использованию
				// Buffer is a variable-sized buffer empty buffer ready to use.
				var buf bytes.Buffer
				buffer := buf.Bytes()

				// Функция копирование через буфер. Like Copy, but only via buffer
				if _, err := io.CopyBuffer(c.Writer, file, buffer); err != nil {
					log.Printf("Error CopyBuffer: %v\n", err)
					c.Status(500) // Ошибка сервиса
				}

				done <- 1
				for range ch1 {
					//nil
				}
			}()
			<-done
			go func() {
				wg.Wait()
				close(done)
				close(ch1)
			}()
			secs := time.Since(start).Seconds()
			log.Printf("%.2fs Time of selector request\n", secs)

		case <-ch2:
			wg.Add(1)
			start := time.Now()
			done := make(chan int, 1)

			go func() {
				defer wg.Done()
				file, err := os.Open("download/" + filename)
				if err != nil {
					c.String(http.StatusNotFound, "Video not found.")
					return
				}
				defer file.Close()

				c.Header("Content-Type", "video/avi")

				var buf2 bytes.Buffer
				buffer2 := buf2.Bytes()

				if _, err := io.CopyBuffer(c.Writer, file, buffer2); err != nil {
					log.Printf("Error CopyBuffer: %v\n", err)
					c.Status(500) // Ошибка сервиса
				}

				done <- 1
				for range ch2 {
					//nil
				}
			}()
			<-done
			go func() {
				wg.Wait()
				close(done)
				close(ch2)
			}()
			secs := time.Since(start).Seconds()
			log.Printf("%.2fs Time of selector request\n", secs)

		case <-ch3:
			wg.Add(1)
			start := time.Now()
			done := make(chan int, 1)

			go func() {
				defer wg.Done()
				file, err := os.Open("download/" + filename)
				if err != nil {
					c.String(http.StatusNotFound, "Video not found.")
					return
				}
				defer file.Close()

				c.Header("Content-Type", "video/mkv")

				var buf3 bytes.Buffer
				buffer3 := buf3.Bytes()

				if _, err := io.CopyBuffer(c.Writer, file, buffer3); err != nil {
					log.Printf("Error CopyBuffer: %v\n", err)
					c.Status(500) // Ошибка сервиса
				}

				done <- 1
				for range ch3 {
					//nil
				}
			}()
			<-done
			go func() {
				wg.Wait()
				close(done)
				close(ch3)
			}()
			secs := time.Since(start).Seconds()
			log.Printf("%.2fs Time of selector request\n", secs)

		case <-ch4:
			wg.Add(1)
			start := time.Now()
			done := make(chan int, 1)

			go func() {
				defer wg.Done()
				file, err := os.Open("download/" + filename)
				if err != nil {
					c.String(http.StatusNotFound, "Video not found.")
					return
				}
				defer file.Close()

				c.Header("Content-Type", "image/gif")

				var buf bytes.Buffer
				buffer := buf.Bytes()

				if _, err := io.CopyBuffer(c.Writer, file, buffer); err != nil {
					log.Printf("Error CopyBuffer: %v\n", err)
					c.Status(500) // Ошибка сервиса
				}

				done <- 1
				for range ch4 {
					//nil
				}
			}()
			<-done
			go func() {
				wg.Wait()
				close(done)
				close(ch4)
			}()
			secs := time.Since(start).Seconds()
			log.Printf("%.2fs Time of selector request\n", secs)

		default:
			log.Println("Неизвестный формат")
			c.Status(500) // Ошибка сервиса
		}
	})

	/////////////////////////////////
	// Uploader. Логика аплоадера
	router.POST("/upstream/:filename", func(c *gin.Context) {

		// Get filename
		filename := c.Param("filename")
		if filename == "" {
			log.Println("Filename is empty!")
			c.Status(400)
		}

		// Get uuid
		uuid := c.GetHeader("user-uuid")
		if uuid == "" {
			log.Println("Uuid is empty!")
			c.Status(400)
		}

		var wg sync.WaitGroup

		var mu sync.Mutex
		chup := make(chan string, 1)
		chid := make(chan string, 1)
		done := make(chan struct{})

		mu.Lock()
		defer mu.Unlock()
		chup <- filename
		chid <- uuid

		wg.Add(1)
		start := time.Now()
		go func() { // Вызов функции создания бакета
			defer wg.Done()

			// Вызов метода создания бакета
			err := makerBucket(chid)
			if err != nil { // TODO обработчик сообщений для фронта
				//c.Status(500) // Ошибка сервиса

				wg.Add(1)
				go func() { // Пишем в существующий бакет со значением uuid
					defer wg.Done()

					c.Header("Content-Type", "application/octet-stream")

					for range <-chid {
						//nil опустошаем канал со строкой сообщения api
					}
					chid <- uuid // Передаем в него uuid

					// Вызов метода загрузчика в бакет
					objname, size := uploader(chup, chid)
					if size == 0 {
						c.Status(500) // Ошибка сервиса, для отладки
					} else {
						log.Printf("Success uploaded %sof size %d\n", objname, size)
					}
					done <- struct{}{}

					for range <-chup {
						//nil
					}
					for range <-chid {
						//nil
					}
				}()
				<-done
				go func() {
					wg.Wait()
					close(chup) //Закрываем канал, собираем горутины
					close(chid)
					close(done)
				}()
			}

			done <- struct{}{}
			for range <-chid {
				//nil
			}
		}()
		<-done
		go func() {
			wg.Wait()
			close(chid) //Закрываем канал, собираем горутины
			close(done)
		}()

		secs := time.Since(start).Seconds()
		log.Printf("%.2fs Time of upload request\n", secs)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Демаршалинг json-файла конфига
func parseConfig() *Params {
	file, err := ioutil.ReadFile("rconfig.json")
	if err != nil {
		log.Fatalf("Error to reading config file: %s", err)
	}

	var pm Params
	if err := json.Unmarshal([]byte(file), &pm); err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err.Error())
	}
	// log.Printf("Params struct: %#v\n", pm)
	return &pm
}
