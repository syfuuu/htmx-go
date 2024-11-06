package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
  templates *template.Template
}

func (t* Templates) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
  return &Templates{
    templates: template.Must(template.ParseGlob("../client/views/*.html")),
  }
}

type Contact struct {
  Name string
  Email string
}

func newContact(name string, email string) Contact {
  return Contact{
    Name: name,
    Email: email,
  }
}

type Contacts = []Contact

type Data struct {
  Contacts Contacts
}

func (d Data) hasEmail(email string) bool {
  for _, contact := range d.Contacts {
    if contact.Email == email {
      return true
    }
  }
  return false
}

func newData() Data {
  return Data{
    Contacts: []Contact{
      newContact("Saif", "saif@gmail.com"),
      newContact("Saif 2", "saif2@gmail.com"),
      newContact("Saif 3", "saif3@gmail.com"),
    },
  }
}

type FormData struct {
  Values map[string]string
  Errors map[string]string
}

func newFormData() FormData {
  return FormData{
    Values: make(map[string]string),
    Errors: make(map[string]string),
  }
}

func main() {
  e := echo.New()
  e.Use(middleware.Logger())

  e.Renderer = newTemplate()

  data := newData()
  e.GET("/", func(c echo.Context) error {
    return c.Render(200, "index", data)
  })

  e.POST("/contacts", func(c echo.Context) error {
    name := c.FormValue("name")
    email := c.FormValue("email")

    if data.hasEmail(email) {
      return c.Render(400, "form", data)
    }

    data.Contacts = append(data.Contacts, newContact(name, email))

    return c.Render(200, "display", data)
  })

  e.Logger.Fatal(e.Start(":8080"))
}
