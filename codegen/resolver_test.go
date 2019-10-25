package codegen_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codegen"
)

var _ = Describe("Resolver", func() {
	var spec *codegen.SpecDescriptor

	Describe("Components", func() {
		Describe("Schemas", func() {
			ItResolvesPrimitiveType := func(name string, index int) {
				message := fmt.Sprintf("resolve the %s type at position %d successfully", name, index)

				It(message, func() {
					typeSpec := spec.Schemas[index]

					Expect(typeSpec.Name).To(Equal(name))
					Expect(typeSpec.IsPrimitive).To(BeTrue())
				})
			}

			ItResolvesEnumType := func(name string, values []interface{}, index int) {
				message := fmt.Sprintf("resolve the %s type at position %d successfully", name, index)

				It(message, func() {
					typeSpec := spec.Schemas[index]

					Expect(typeSpec.Name).To(Equal(name))
					Expect(typeSpec.IsEnum).To(BeTrue())
					Expect(typeSpec.Metadata).To(HaveKeyWithValue("enum", values))
				})
			}

			Describe("Integer", func() {
				BeforeEach(func() {
					spec = resolve("schemas-integer.yaml")
					Expect(spec.Schemas).To(HaveLen(4))
				})

				Describe("int32", func() {
					ItResolvesPrimitiveType("int32", 0)
					ItResolvesPrimitiveType("int32", 1)
				})

				Describe("int64", func() {
					ItResolvesPrimitiveType("int64", 2)
					ItResolvesPrimitiveType("int64", 3)
				})
			})

			Describe("Number", func() {
				BeforeEach(func() {
					spec = resolve("schemas-number.yaml")
					Expect(spec.Schemas).To(HaveLen(4))
				})

				Describe("double", func() {
					ItResolvesPrimitiveType("double", 0)
					ItResolvesPrimitiveType("double", 1)
				})

				Describe("float", func() {
					ItResolvesPrimitiveType("float", 2)
					ItResolvesPrimitiveType("float", 3)
				})
			})

			Describe("String", func() {
				BeforeEach(func() {
					spec = resolve("schemas-string.yaml")
					Expect(spec.Schemas).To(HaveLen(12))
				})

				Describe("binary", func() {
					ItResolvesPrimitiveType("binary", 0)
					ItResolvesPrimitiveType("binary", 1)
				})

				Describe("byte", func() {
					ItResolvesPrimitiveType("byte", 2)
					ItResolvesPrimitiveType("byte", 3)
				})

				Describe("date", func() {
					ItResolvesPrimitiveType("date", 4)
					ItResolvesPrimitiveType("date", 5)
				})

				Describe("date-time", func() {
					ItResolvesPrimitiveType("date-time", 6)
					ItResolvesPrimitiveType("date-time", 7)
				})

				Describe("string", func() {
					ItResolvesPrimitiveType("string", 8)
					ItResolvesPrimitiveType("string", 9)
				})

				Describe("uuid", func() {
					ItResolvesPrimitiveType("uuid", 10)
					ItResolvesPrimitiveType("uuid", 11)
				})
			})

			Describe("Enum", func() {
				values := []interface{}{"pending", "completed"}

				BeforeEach(func() {
					spec = resolve("schemas-enum.yaml")
					Expect(spec.Schemas).To(HaveLen(2))
				})

				ItResolvesEnumType("AccountStatus", values, 0)
				ItResolvesEnumType("TransactionStatus", values, 1)
			})

			Describe("Array", func() {
			})
		})

		Describe("Parameters", func() {
		})

		Describe("Responses", func() {
		})

		Describe("Requests", func() {
		})
	})
})
