package golang_test

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
	"github.com/phogolabs/stride/plugin/golang"
)

var _ = Describe("SchemaGenerator", func() {
	var generator *golang.SchemaGenerator

	BeforeEach(func() {
		generator = &golang.SchemaGenerator{
			Path: tmpdir(),
		}

		Expect(os.MkdirAll(generator.Path, 0755)).To(Succeed())
	})

	Context("when the descriptor is alias", func() {
		BeforeEach(func() {
			descriptor := &codegen.TypeDescriptor{
				Name:    "ID",
				IsAlias: true,
				Element: &codegen.TypeDescriptor{
					Name:        "string",
					IsPrimitive: true,
				},
			}

			generator.Collection = append(generator.Collection, descriptor)
		})

		It("generates the schema successfully", func() {
			file := generator.Generate()
			Expect(file).NotTo(BeNil())

			path := file.Name()
			Expect(filepath.Base(path)).To(Equal("schema.go"))

			buffer := &bytes.Buffer{}
			_, err := file.WriteTo(buffer)
			Expect(err).To(BeNil())

			var (
				scanner = bufio.NewScanner(buffer)
				line    = 0
			)

			for scanner.Scan() {
				text := scanner.Text()

				switch line {
				case 0:
					Expect(text).To(Equal("package service"))
				case 2:
					Expect(text).To(Equal("// ID is a type auto-generated from OpenAPI spec"))
				case 3:
					Expect(text).To(Equal("// stride:generate id"))
				case 4:
					Expect(text).To(Equal("type ID string"))
				}

				line = line + 1
			}
		})
	})

	Context("when the descriptor is array", func() {
		BeforeEach(func() {
			descriptor := &codegen.TypeDescriptor{
				Name:    "Names",
				IsArray: true,
				Element: &codegen.TypeDescriptor{
					Name:        "string",
					IsPrimitive: true,
				},
			}

			generator.Collection = append(generator.Collection, descriptor)
		})

		It("generates the schema successfully", func() {
			file := generator.Generate()
			Expect(file).NotTo(BeNil())

			path := file.Name()
			Expect(filepath.Base(path)).To(Equal("schema.go"))

			buffer := &bytes.Buffer{}
			_, err := file.WriteTo(buffer)
			Expect(err).To(BeNil())

			var (
				scanner = bufio.NewScanner(buffer)
				line    = 0
			)

			for scanner.Scan() {
				text := scanner.Text()

				switch line {
				case 0:
					Expect(text).To(Equal("package service"))
				case 2:
					Expect(text).To(Equal("// Names is a type auto-generated from OpenAPI spec"))
				case 3:
					Expect(text).To(Equal("// stride:generate names"))
				case 4:
					Expect(text).To(Equal("type Names []string"))
				}

				line = line + 1
			}
		})
	})

	Context("when the descriptor is class", func() {
		BeforeEach(func() {
			descriptor := &codegen.TypeDescriptor{
				Name:    "User",
				IsClass: true,
				Properties: codegen.PropertyDescriptorCollection{
					&codegen.PropertyDescriptor{
						Name: "ID",
						PropertyType: &codegen.TypeDescriptor{
							Name:        "string",
							IsPrimitive: true,
						},
					},
				},
			}

			generator.Collection = append(generator.Collection, descriptor)
		})

		It("generates the schema successfully", func() {
			file := generator.Generate()
			Expect(file).NotTo(BeNil())

			path := file.Name()
			Expect(filepath.Base(path)).To(Equal("schema.go"))

			buffer := &bytes.Buffer{}
			_, err := file.WriteTo(buffer)
			Expect(err).To(BeNil())

			var (
				scanner = bufio.NewScanner(buffer)
				line    = 0
			)

			for scanner.Scan() {
				text := scanner.Text()

				switch line {
				case 0:
					Expect(text).To(Equal("package service"))
				case 2:
					Expect(text).To(Equal("// User is a type auto-generated from OpenAPI spec"))
				case 3:
					Expect(text).To(Equal("// stride:generate user"))
				case 4:
					Expect(text).To(Equal("type User struct {"))
				case 5:
					Expect(text).To(Equal("\t// stride:generate id"))
				case 6:
					Expect(text).To(Equal("\tID string `json:\"ID,omitempty\" validate:\"-\"`"))
				case 7:
					Expect(text).To(Equal("}"))
				}

				line = line + 1
			}
		})
	})

	Context("when the descriptor is enum", func() {
		BeforeEach(func() {
			descriptor := &codegen.TypeDescriptor{
				Name:   "Status",
				IsEnum: true,
				Metadata: codegen.Metadata{
					"values": []interface{}{"pending", "running", "completed"},
				},
			}

			generator.Collection = append(generator.Collection, descriptor)
		})

		It("generates the schema successfully", func() {
			file := generator.Generate()
			Expect(file).NotTo(BeNil())

			path := file.Name()
			Expect(filepath.Base(path)).To(Equal("schema.go"))

			buffer := &bytes.Buffer{}
			_, err := file.WriteTo(buffer)
			Expect(err).To(BeNil())

			var (
				scanner = bufio.NewScanner(buffer)
				line    = 0
			)

			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())

				switch line {
				case 0:
					Expect(text).To(Equal("package service"))
				case 2:
					Expect(text).To(Equal("// Status is a type auto-generated from OpenAPI spec"))
				case 3:
					Expect(text).To(Equal("// stride:generate status"))
				case 4:
					Expect(text).To(Equal("type Status string"))
				case 7:
					Expect(text).To(Equal("// StatusPending is a \"pending\" constant auto-generated from OpenAPI spec"))
				case 8:
					Expect(text).To(Equal("// stride:generate status-pending"))
				case 9:
					Expect(text).To(Equal("StatusPending Status = \"pending\""))
				}

				line = line + 1
			}
		})
	})
})
