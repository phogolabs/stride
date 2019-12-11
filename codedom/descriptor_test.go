package codedom_test

import (
	"fmt"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/stride/codedom"
)

var _ = Describe("TypeDescriptor", func() {
	Describe("Kind", func() {
		ItReturnTheKind := func(name, expected string, primitive, nullable bool) {
			It(fmt.Sprintf("returns %s successfully", expected), func() {
				descriptor := &codedom.TypeDescriptor{
					Name:        name,
					IsPrimitive: primitive,
					IsNullable:  nullable,
				}

				Expect(descriptor.Kind()).To(Equal(expected))
			})
		}

		ItReturnTheKind("date-time", "time.Time", true, false)
		ItReturnTheKind("date", "time.Time", true, false)
		ItReturnTheKind("uuid", "schema.UUID", true, false)
		ItReturnTheKind("int", "int", true, false)
		ItReturnTheKind("int", "*int", true, true)
	})

	Describe("Tags", func() {
		It("returns a tag collection successfully", func() {
			float64Ptr := func(v float64) *float64 {
				return &v
			}

			kind := &codedom.TypeDescriptor{
				Default: "99.9",
				Metadata: codedom.Metadata{
					"unique":        true,
					"min":           float64Ptr(10.0),
					"min_exclusive": true,
					"max":           float64Ptr(20.0),
					"max_exclusive": true,
					"multiple_of":   float64Ptr(2),
					"pattern":       "[a-Z]",
					"values":        []interface{}{1, 2, 3},
				},
			}

			tags := kind.Tags(true)
			Expect(tags).To(HaveLen(2))
			Expect(tags[0].Key).To(Equal("validate"))
			Expect(tags[0].Name).To(Equal("required"))
			Expect(tags[0].Options).To(ContainElement("unique"))
			Expect(tags[0].Options).To(ContainElement("oneof=1 2 3"))
			Expect(tags[0].Options).To(ContainElement("gt=10"))
			Expect(tags[0].Options).To(ContainElement("lt=20"))

			Expect(tags[1].Key).To(Equal("default"))
			Expect(tags[1].Name).To(Equal("99.9"))
		})

		Context("when the exlusive is disabled", func() {
			It("returns a tag collection successfully", func() {
				float64Ptr := func(v float64) *float64 {
					return &v
				}

				kind := &codedom.TypeDescriptor{
					Default: "99.9",
					Metadata: codedom.Metadata{
						"unique":        true,
						"min":           float64Ptr(10.0),
						"min_exclusive": false,
						"max":           float64Ptr(20.0),
						"max_exclusive": false,
						"multiple_of":   float64Ptr(2),
						"pattern":       "[a-Z]",
						"values":        []interface{}{1, 2, 3},
					},
				}

				tags := kind.Tags(true)
				Expect(tags).To(HaveLen(2))
				Expect(tags[0].Key).To(Equal("validate"))
				Expect(tags[0].Name).To(Equal("required"))
				Expect(tags[0].Options).To(ContainElement("unique"))
				Expect(tags[0].Options).To(ContainElement("oneof=1 2 3"))
				Expect(tags[0].Options).To(ContainElement("gt=10"))
				Expect(tags[0].Options).To(ContainElement("lt=20"))

				Expect(tags[1].Key).To(Equal("default"))
				Expect(tags[1].Name).To(Equal("99.9"))

				Expect(kind.HasProperties()).To(BeFalse())
			})
		})
	})
})

var _ = Describe("TypeDescriptorMap", func() {
	var descriptor *codedom.TypeDescriptor

	BeforeEach(func() {
		descriptor = &codedom.TypeDescriptor{
			Name: "my-type",
		}
	})

	Describe("Add", func() {
		It("adds a descriptor successfully", func() {
			kv := codedom.TypeDescriptorMap{}
			kv.Add(descriptor)
			Expect(kv).To(HaveKeyWithValue("my-type", descriptor))
		})
	})

	Describe("Get", func() {
		It("gets a descriptor successfully", func() {
			kv := codedom.TypeDescriptorMap{}
			kv.Add(descriptor)
			Expect(kv).To(HaveLen(1))
			Expect(kv.Get("my-type")).To(Equal(descriptor))
		})
	})

	Describe("Clear", func() {
		It("clears a descriptor map successfully", func() {
			kv := codedom.TypeDescriptorMap{}
			kv.Add(descriptor)

			Expect(kv).To(HaveLen(1))
			kv.Clear()
			Expect(kv).To(HaveLen(0))
		})
	})

	Describe("Collection", func() {
		It("returns a descriptor collection successfully", func() {
			kv := codedom.TypeDescriptorMap{}
			kv.Add(descriptor)
			Expect(kv.Collection()).To(ContainElement(descriptor))
		})
	})
})

var _ = Describe("TypeDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.TypeDescriptorCollection{}
			items = append(items, &codedom.TypeDescriptor{Name: "BankID"})
			items = append(items, &codedom.TypeDescriptor{Name: "ID"})
			items = append(items, &codedom.TypeDescriptor{Name: "Address"})

			sort.Sort(items)

			Expect(items[0].Name).To(Equal("Address"))
			Expect(items[1].Name).To(Equal("BankID"))
			Expect(items[2].Name).To(Equal("ID"))
		})
	})
})

var _ = Describe("PropertyDescriptor", func() {
	Describe("Tags", func() {
		It("returns a tag collection successfully", func() {
			property := &codedom.PropertyDescriptor{
				Name:     "bank-id",
				Required: false,
				PropertyType: &codedom.TypeDescriptor{
					Name:    "string",
					Default: "hello",
				},
			}

			tags := property.Tags()
			Expect(tags).To(HaveLen(6))
			Expect(tags[0].Key).To(Equal("json"))
			Expect(tags[0].Name).To(Equal("bank-id"))
			Expect(tags[0].Options).To(ContainElement("omitempty"))

			Expect(tags[1].Key).To(Equal("xml"))
			Expect(tags[1].Name).To(Equal("bank-id"))
			Expect(tags[1].Options).To(ContainElement("omitempty"))

			Expect(tags[2].Key).To(Equal("form"))
			Expect(tags[2].Name).To(Equal("bank-id"))
			Expect(tags[2].Options).To(ContainElement("omitempty"))

			Expect(tags[3].Key).To(Equal("field"))
			Expect(tags[3].Name).To(Equal("bank-id"))
			Expect(tags[3].Options).To(ContainElement("omitempty"))

			Expect(tags[4].Key).To(Equal("validate"))
			Expect(tags[4].Name).To(Equal("-"))
			Expect(tags[4].Options).To(BeEmpty())

			Expect(tags[5].Key).To(Equal("default"))
			Expect(tags[5].Name).To(Equal("hello"))
			Expect(tags[5].Options).To(BeEmpty())
		})
	})
})

var _ = Describe("PropertyDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.PropertyDescriptorCollection{}
			items = append(items, &codedom.PropertyDescriptor{Name: "bank-id"})
			items = append(items, &codedom.PropertyDescriptor{Name: "address"})

			sort.Sort(items)

			Expect(items[0].Name).To(Equal("address"))
			Expect(items[1].Name).To(Equal("bank-id"))
		})
	})
})

var _ = Describe("ParameterDescriptor", func() {
	Describe("Tags", func() {
		It("returns a tag collection successfully", func() {
			param := &codedom.ParameterDescriptor{
				Name:     "bank-id",
				Style:    "matrix",
				In:       "header",
				Explode:  true,
				Required: true,
				ParameterType: &codedom.TypeDescriptor{
					Name: "string",
				},
			}

			tags := param.Tags()
			Expect(tags).To(HaveLen(2))
			Expect(tags[0].Key).To(Equal("header"))
			Expect(tags[0].Name).To(Equal("bank-id"))
			Expect(tags[0].Options).To(ContainElement("matrix"))

			Expect(tags[1].Key).To(Equal("validate"))
			Expect(tags[1].Name).To(Equal("required"))
			Expect(tags[1].Options).To(BeEmpty())
		})
	})
})

var _ = Describe("ParameterDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.ParameterDescriptorCollection{}
			items = append(items, &codedom.ParameterDescriptor{Name: "bank-id"})
			items = append(items, &codedom.ParameterDescriptor{Name: "address"})

			sort.Sort(items)

			Expect(items[0].Name).To(Equal("address"))
			Expect(items[1].Name).To(Equal("bank-id"))
		})
	})
})

var _ = Describe("RequestDescriptor", func() {
	//TODO: implement it if your need
})

var _ = Describe("RequestDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.RequestDescriptorCollection{}
			items = append(items, &codedom.RequestDescriptor{ContentType: "application/xml"})
			items = append(items, &codedom.RequestDescriptor{ContentType: "application/json"})

			sort.Sort(items)

			Expect(items[0].ContentType).To(Equal("application/json"))
			Expect(items[1].ContentType).To(Equal("application/xml"))
		})
	})
})

var _ = Describe("ResponseDescriptor", func() {
	//TODO: implement it if your need
})

var _ = Describe("ResponseDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.ResponseDescriptorCollection{}
			items = append(items, &codedom.ResponseDescriptor{ContentType: "application/xml", Code: 200})
			items = append(items, &codedom.ResponseDescriptor{ContentType: "application/json", Code: 200})
			items = append(items, &codedom.ResponseDescriptor{ContentType: "application/xml", Code: 201})
			items = append(items, &codedom.ResponseDescriptor{ContentType: "application/json", Code: 201})

			sort.Sort(items)

			Expect(items[0].ContentType).To(Equal("application/json"))
			Expect(items[0].Code).To(Equal(200))
			Expect(items[1].ContentType).To(Equal("application/json"))
			Expect(items[1].Code).To(Equal(201))
			Expect(items[2].ContentType).To(Equal("application/xml"))
			Expect(items[2].Code).To(Equal(200))
			Expect(items[3].ContentType).To(Equal("application/xml"))
			Expect(items[3].Code).To(Equal(201))
		})
	})
})

var _ = Describe("ControllerDescriptor", func() {
	//TODO: implement it if your need
})

var _ = Describe("ControllerDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.ControllerDescriptorCollection{}
			items = append(items, &codedom.ControllerDescriptor{Name: "account-api"})
			items = append(items, &codedom.ControllerDescriptor{Name: "user-api"})
			items = append(items, &codedom.ControllerDescriptor{Name: "card-api"})

			sort.Sort(items)

			Expect(items[0].Name).To(Equal("account-api"))
			Expect(items[1].Name).To(Equal("card-api"))
			Expect(items[2].Name).To(Equal("user-api"))
		})
	})
})

var _ = Describe("ControllerDescriptorMap", func() {
	var descriptor *codedom.ControllerDescriptor

	BeforeEach(func() {
		descriptor = &codedom.ControllerDescriptor{
			Name: "user-api",
		}
	})

	Describe("Add", func() {
		It("adds a descriptor successfully", func() {
			kv := codedom.ControllerDescriptorMap{}
			kv.Add(descriptor)
			kv.Add(descriptor)
			Expect(kv).To(HaveKeyWithValue("user-api", descriptor))
		})
	})

	Describe("Get", func() {
		It("gets a descriptor successfully", func() {
			kv := codedom.ControllerDescriptorMap{}
			kv.Add(descriptor)
			Expect(kv).To(HaveLen(1))
			Expect(kv.Get("user-api")).To(Equal(descriptor))
		})

		Context("when the descriptor does not exist", func() {
			It("gets a descriptor successfully", func() {
				kv := codedom.ControllerDescriptorMap{}
				descriptor = kv.Get("user-api")
				Expect(descriptor).NotTo(BeNil())
				Expect(descriptor.Name).To(Equal("user-api"))
			})
		})
	})

	Describe("Clear", func() {
		It("clears a descriptor map successfully", func() {
			kv := codedom.ControllerDescriptorMap{}
			kv.Add(descriptor)

			Expect(kv).To(HaveLen(1))
			kv.Clear()
			Expect(kv).To(HaveLen(0))
		})
	})

	Describe("Collection", func() {
		It("returns a descriptor collection successfully", func() {
			kv := codedom.ControllerDescriptorMap{}
			kv.Add(descriptor)
			Expect(kv.Collection()).To(ContainElement(descriptor))
		})
	})
})

var _ = Describe("OperationDescriptor", func() {
	//TODO: implement it if your need
})

var _ = Describe("OperationDescriptorCollection", func() {
	Describe("Sort", func() {
		It("sorts the items successfully", func() {
			items := codedom.OperationDescriptorCollection{}
			items = append(items, &codedom.OperationDescriptor{Name: "create"})
			items = append(items, &codedom.OperationDescriptor{Name: "update"})
			items = append(items, &codedom.OperationDescriptor{Name: "delete"})

			sort.Sort(items)

			Expect(items[0].Name).To(Equal("create"))
			Expect(items[1].Name).To(Equal("delete"))
			Expect(items[2].Name).To(Equal("update"))
		})
	})
})

var _ = Describe("TagDescriptor", func() {
	//TODO: implement it if your need
})

var _ = Describe("TagDescriptorCollection", func() {
	Describe("String", func() {
		It("returns the tag string successfully", func() {
			items := codedom.TagDescriptorCollection{}
			items = append(items, &codedom.TagDescriptor{Key: "json", Name: "name"})
			items = append(items, &codedom.TagDescriptor{Key: "xml", Name: "name"})

			Expect(items.String()).To(Equal("`json:\"name\" xml:\"name\"`"))
		})

		Context("when the collection is empty", func() {
			It("returns the tag string successfully", func() {
				items := codedom.TagDescriptorCollection{}
				Expect(items.String()).To(BeEmpty())
			})
		})
	})
})
