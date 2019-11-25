package syntax_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dave/dst"
	"github.com/phogolabs/stride/fake"
	"github.com/phogolabs/stride/plugin/syntax"
)

var _ = Describe("File", func() {
	Describe("OpenFile", func() {
		It("opens the file successfully", func() {
			file, err := syntax.OpenFile("../../fixture/code/user_defined_source.go.fixture")
			Expect(err).To(BeNil())
			Expect(file).NotTo(BeNil())
			Expect(file.Name()).To(Equal("../../fixture/code/user_defined_source.go.fixture"))
			Expect(file.Node()).NotTo(BeNil())
		})

		Context("when the file does not exist", func() {
			It("returns an error", func() {
				file, err := syntax.OpenFile("./file-not-exist.go")
				Expect(err).To(MatchError("open ./file-not-exist.go: no such file or directory"))
				Expect(file).To(BeNil())
			})
		})
	})

	Describe("WriteTo", func() {
		It("writes the file to a writer", func() {
			file, err := syntax.OpenFile("../../fixture/code/user_defined_source.go.fixture")
			Expect(err).To(BeNil())
			Expect(file).NotTo(BeNil())

			buffer := &bytes.Buffer{}

			_, err = file.WriteTo(buffer)
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(ContainSubstring("type User struct"))
		})

		Context("when the writer fails", func() {
			It("returns an error", func() {
				file, err := syntax.OpenFile("../../fixture/code/user_defined_source.go.fixture")
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
			file := syntax.NewFile(tmpfile())
			Expect(file.Sync()).To(Succeed())
			Expect(file.Name()).To(BeAnExistingFile())
		})

		Context("when the file cannot be written", func() {
			It("returns an error", func() {
				file := syntax.NewFile("./unknown/root.go")
				Expect(file.Sync()).To(MatchError("open ./unknown/root.go: no such file or directory"))
			})
		})
	})

	Describe("Literal", func() {
		It("returns a literal", func() {
			file := syntax.NewFile("model.go")
			Expect(file.Literal("ID")).NotTo(BeNil())
		})
	})

	Describe("Array", func() {
		It("returns a array", func() {
			file := syntax.NewFile("model.go")
			Expect(file.Array("ID")).NotTo(BeNil())
		})
	})

	Describe("Struct", func() {
		It("returns a struct", func() {
			file := syntax.NewFile("model.go")
			Expect(file.Struct("ID")).NotTo(BeNil())
		})
	})
})

var _ = Describe("Literal", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := syntax.NewLiteralType("User")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := syntax.NewLiteralType("User")
			Expect(spec.Name()).To(Equal("User"))
		})
	})

	Describe("Element", func() {
		It("sets the element", func() {
			spec := syntax.NewLiteralType("ID")
			spec.Element("string")
			Expect(spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.Ident).Name).To(Equal("string"))
		})
	})
})

var _ = Describe("Array", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := syntax.NewArrayType("Names")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := syntax.NewArrayType("Names")
			Expect(spec.Name()).To(Equal("Names"))
		})
	})

	Describe("Element", func() {
		It("sets the element", func() {
			spec := syntax.NewArrayType("Names")
			spec.Element("string")
			Expect(spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.ArrayType).Elt.(*dst.Ident).Name).To(Equal("string"))
		})
	})
})

var _ = Describe("Struct", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := syntax.NewStructType("User")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := syntax.NewStructType("User")
			Expect(spec.Name()).To(Equal("User"))
		})
	})

	Describe("AddField", func() {
		It("adds a new field", func() {
			spec := syntax.NewStructType("User")
			spec.AddField("ID", "string")

			structType := spec.Node().Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
			Expect(structType.Fields.List).To(HaveLen(1))
			Expect(structType.Fields.List[0].Names[0].Name).To(Equal("ID"))
		})
	})

	Describe("Function", func() {
		It("returns a function", func() {
			spec := syntax.NewStructType("User")
			Expect(spec.Function("Status")).NotTo(BeNil())
		})
	})
})

var _ = Describe("Function", func() {
	Describe("Commentf", func() {
		It("comments the type", func() {
			spec := syntax.NewFunctionType("AddUser")
			spec.Commentf("my comment")
			Expect(spec.Node().Decs.Start.All()).To(ContainElement("// my comment"))
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			spec := syntax.NewFunctionType("AddUser")
			Expect(spec.Name()).To(Equal("AddUser"))
		})
	})

	Describe("AddReceiver", func() {
		It("adds the receiver", func() {
			spec := syntax.NewFunctionType("AddUser")
			spec.AddReceiver("x", "Controller")
			Expect(spec.Node().Recv.List).To(HaveLen(1))
			Expect(spec.Node().Recv.List[0].Names[0].Name).To(Equal("x"))
			Expect(spec.Node().Recv.List[0].Type.(*dst.Ident).Name).To(Equal("Controller"))
		})
	})

	Describe("AddParam", func() {
		It("adds the param", func() {
			spec := syntax.NewFunctionType("AddUser")
			spec.AddParam("name", "string")
			Expect(spec.Node().Type.Params.List).To(HaveLen(1))
			Expect(spec.Node().Type.Params.List[0].Names[0].Name).To(Equal("name"))
			Expect(spec.Node().Type.Params.List[0].Type.(*dst.Ident).Name).To(Equal("string"))
		})
	})

	Describe("AddReturn", func() {
		It("adds the return param", func() {
			spec := syntax.NewFunctionType("Status")
			spec.AddReturn("int")
			Expect(spec.Node().Type.Results.List).To(HaveLen(1))
			Expect(spec.Node().Type.Results.List[0].Type.(*dst.Ident).Name).To(Equal("int"))
		})
	})

	Describe("Body", func() {
		It("returns the body", func() {
			spec := syntax.NewFunctionType("Status")
			Expect(spec.Body()).NotTo(BeNil())
		})
	})
})

var _ = Describe("Block", func() {
	Describe("Append", func() {
		It("writes text to a block", func() {
			block := syntax.NewBlockType()
			block.Append("fmt.Println(123)")
			Expect(block.Build()).To(Succeed())
			Expect(block.Node().List).To(HaveLen(1))
		})
	})

	Describe("AppendComment", func() {
		It("writes comments to a block", func() {
			block := syntax.NewBlockType()
			block.AppendComment()
			Expect(block.Build()).To(Succeed())
		})
	})

	Describe("Build", func() {
		It("builds the block", func() {
			block := syntax.NewBlockType()
			block.Append("fmt.Println(123)")
			Expect(block.Build()).To(Succeed())
			Expect(block.Node().List).To(HaveLen(1))
		})
	})
})
