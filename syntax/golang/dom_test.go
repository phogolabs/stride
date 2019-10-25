package golang_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dave/dst"
	"github.com/phogolabs/stride/fake"
	"github.com/phogolabs/stride/syntax/golang"
)

var _ = Describe("File", func() {
	Describe("OpenFile", func() {
		It("opens the file successfully", func() {
			file, err := golang.OpenFile("../../fixture/code/user_defined_source.go.fixture")
			Expect(err).To(BeNil())
			Expect(file).NotTo(BeNil())
			Expect(file.Name()).To(Equal("../../fixture/code/user_defined_source.go.fixture"))
			Expect(file.Node()).NotTo(BeNil())
		})

		Context("when the file does not exist", func() {
			It("returns an error", func() {
				file, err := golang.OpenFile("./file-not-exist.go")
				Expect(err).To(MatchError("open ./file-not-exist.go: no such file or directory"))
				Expect(file).To(BeNil())
			})
		})
	})

	Describe("WriteTo", func() {
		It("writes the file to a writer", func() {
			file, err := golang.OpenFile("../../fixture/code/user_defined_source.go.fixture")
			Expect(err).To(BeNil())
			Expect(file).NotTo(BeNil())

			buffer := &bytes.Buffer{}

			_, err = file.WriteTo(buffer)
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(ContainSubstring("type User struct"))
		})

		Context("when the writer fails", func() {
			It("returns an error", func() {
				file, err := golang.OpenFile("../../fixture/code/user_defined_source.go.fixture")
				Expect(err).To(BeNil())
				Expect(file).NotTo(BeNil())

				writer := &fake.Writer{}
				writer.WriteReturns(0, fmt.Errorf("oh no"))

				_, err = file.WriteTo(writer)
				Expect(err).To(MatchError("oh no"))
			})
		})
	})

	Describe("Sync", func() {
		It("writes the file to the disk", func() {
			file := golang.NewFile(tmpfile())
			Expect(file.Sync()).To(Succeed())
			Expect(file.Name()).To(BeAnExistingFile())
		})

		Context("when the file cannot be written", func() {
			It("returns an error", func() {
				file := golang.NewFile("./unknown/root.go")
				Expect(file.Sync()).To(MatchError("open ./unknown/root.go: no such file or directory"))
			})
		})
	})
})

var _ = Describe("Literal", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := golang.NewLiteralType("User")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := golang.NewLiteralType("User")
			Expect(spec.Name()).To(Equal("User"))
		})
	})

	Describe("Element", func() {
		It("sets the element", func() {
			spec := golang.NewLiteralType("ID")
			spec.Element("string")
			Expect(spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.Ident).Name).To(Equal("string"))
		})
	})
})

var _ = Describe("Array", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := golang.NewArrayType("Names")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := golang.NewArrayType("Names")
			Expect(spec.Name()).To(Equal("Names"))
		})
	})

	Describe("Element", func() {
		It("sets the element", func() {
			spec := golang.NewArrayType("Names")
			spec.Element("string")
			Expect(spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.ArrayType).Elt.(*dst.Ident).Name).To(Equal("string"))
		})
	})
})

var _ = Describe("Struct", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := golang.NewStructType("User")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := golang.NewStructType("User")
			Expect(spec.Name()).To(Equal("User"))
		})
	})

	Describe("AddField", func() {
		It("adds a new field", func() {
			spec := golang.NewStructType("User")
			spec.AddField("ID", "string")

			structType := spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
			Expect(structType.Fields.List).To(HaveLen(1))
			Expect(structType.Fields.List[0].Names[0].Name).To(Equal("ID"))
		})
	})
})
