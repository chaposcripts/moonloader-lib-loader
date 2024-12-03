package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	runtime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "list:request", func(a ...interface{}) {

	})
}

func (a *App) RequestLibsData() error {
	fmt.Println("list:request")
	err := loadData()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Ошибка",
			Message: fmt.Sprintf("Ошибка: %s", err.Error()),
		})
		return err
	}
	runtime.EventsEmit(a.ctx, "list:update", true, libsData)
	return nil
}

func (a *App) RequestFolderSelection() error {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "GTA:SA Folder",
	})
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Ошибка",
			Message: fmt.Sprintf("Ошибка: %s", err.Error()),
		})
		return err
	}
	if len(path) > 0 {
		if _, err := os.Stat(path + "\\MoonLoader.asi"); err != nil {
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "Ошибка",
				Message: fmt.Sprintf("Ошибка, файл MoonLoader.asi не найден!\n%s", err.Error()),
			})
			return err
		}
		runtime.EventsEmit(a.ctx, "path:selected", path)
	}
	return nil
}

func (a *App) pushLog(s string) {
	runtime.EventsEmit(a.ctx, "log:push", s)
}

func (a *App) InstallSelectedLibs(path string, selectedLibs []string) error {
	if len(path) == 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{Title: "Ошибка", Message: "Ошибка, укажите папку с игрой"})
		return fmt.Errorf("EMPTY_FOLDER")
	}
	if !dataReceived {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{Title: "Ошибка", Message: "Ошибка, данные не были получены"})
		return fmt.Errorf("DATA_NOT_LOADED")
	}

	// download zip
	runtime.EventsEmit(a.ctx, "status:set", "downloading")
	err := loadZip(zipURL, path+"\\moonloader-libs-temp.zip")
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{Title: "Ошибка", Message: fmt.Sprintf("Ошибка загрузки: %s", err.Error())})
		return err
	}

	// extract zip to temp folder
	runtime.EventsEmit(a.ctx, "status:set", "extracting")
	list, err := extractAllFiles(selectedLibs, path+"\\moonloader-libs-temp.zip", path+"\\moonloader")
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{Title: "Ошибка", Message: fmt.Sprintf("Ошибка распаковки: %s", err.Error())})
		return err
	}
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Успешно!",
		Message: fmt.Sprintf("В ходе установки были установлены следующие файлы: %s", strings.Join((list), "\n")),
	})

	runtime.EventsEmit(a.ctx, "status:set", "copying")
	// move selected libs to game folder
	runtime.EventsEmit(a.ctx, "status:set", "deleting")
	// delete temp folder
	runtime.EventsEmit(a.ctx, "status:set", "none")

	return nil
}
