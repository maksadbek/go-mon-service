// @author: Maksadbek
// @email: a.maksadbek@gmail.com

// пакет для кофигурационных данный
// формат томл используется как конф. файл
package conf

import (
	"io"

	"github.com/BurntSushi/toml"
)

type ErrorStr struct {
	Msg string
}

// структура для конф. данных
type Datastore struct {
	Mysql struct {
		DSN      string
		Interval int
	}
	// анонимная структура для конфигурации для редис
	Redis struct {
		Host    string // хост для редис сервера
		FPrefix string // флит префикс, так будет сохранятся в редис. Например: fleet_202, flit_202, ...
		TPrefix string // трекер префикс, так же как флит префикс. Например: tracker_512341, another_2131
		UPrefix string
	}
}

// структура конф. для сервера
type Server struct {
	IP   string // адрес
	Port string // порт
}

// главный структура для конф.
type App struct {
	DS       Datastore           // база данных
	SRV      Server              // сервер
	Log      Log                 // логирования
	ErrorMsg map[string]ErrorStr `toml:"errors"`
}

// структура для логирования
type Log struct {
	Path string // пут для лог файла
}

// читает данные из конф файла и возврашаеть структуру
func Read(r io.Reader) (config App, err error) {
	_, err = toml.DecodeReader(r, &config)
	return config, err
}
