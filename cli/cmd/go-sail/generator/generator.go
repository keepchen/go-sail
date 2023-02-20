package generator

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Generator interface {
	Generate([]File) error
}

type generator struct {
	opts Options
}

type File struct {
	Path     string
	Template string
}

func (g *generator) Generate(files []File) error {
	for _, file := range files {
		log.Println("-----", file.Path)
		fp := filepath.Join(g.opts.WorkDir, file.Path)
		dir := filepath.Dir(fp)
		//是创建目录
		if file.Template == "" {
			dir = fp
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		if file.Template == "" {
			continue
		}

		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		t, err := template.New(fp).Parse(file.Template)
		if err != nil {
			return err
		}

		err = t.Execute(f, g.opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) GetOptions() Options {
	return g.opts
}

// New returns a new generator struct.
func New(opts Options) Generator {
	return &generator{
		opts: opts,
	}
}
