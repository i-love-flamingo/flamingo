package pugast

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pugast Rendering", func() {
	var p = NewPugAst("/")

	It("Should Render Tags with intendation", func() {
		var testast = p.ParseJson([]byte(`
{
  "type": "Block",
  "nodes": [
    {
      "type": "Doctype",
      "val": "html",
      "line": 1,
      "filename": "layouts/blank.pug"
    },
    {
      "type": "Tag",
      "name": "html",
      "selfClosing": false,
      "block": {
        "type": "Block",
        "nodes": [

          {
            "type": "Tag",
            "name": "body",
            "selfClosing": false,
            "block": {
              "type": "Block",
              "nodes": [
                {
                  "type": "NamedBlock",
                  "nodes": [
                    {
                      "type": "Code",
                      "val": "Site = get('site')",
                      "buffer": false,
                      "mustEscape": false,
                      "isInline": false,
                      "line": 4,
                      "filename": "layouts/default.pug"
                    },
                    {
                      "type": "Tag",
                      "name": "a",
                      "selfClosing": false,
                      "block": {
                        "type": "Block",
                        "nodes": [],
                        "line": 5,
                        "filename": "layouts/default.pug"
                      },
                      "attrs": [
                        {
                          "name": "name",
                          "val": "\"top\"",
                          "mustEscape": true
                        }
                      ],
                      "attributeBlocks": [],
                      "isInline": true,
                      "line": 5,
                      "filename": "layouts/default.pug"
                    },
                    {
                      "type": "Tag",
                      "name": "section",
                      "selfClosing": false,
                      "block": {
                        "type": "Block",
                        "nodes": [
                          {
                            "type": "Tag",
                            "name": "div",
                            "selfClosing": false,
                            "block": {
                              "type": "Block",
                              "nodes": [
                                {
                                  "type": "NamedBlock",
                                  "nodes": [
                                    {
                                      "type": "Tag",
                                      "name": "div",
                                      "selfClosing": false,
                                      "block": {
                                        "type": "Block",
                                        "nodes": [
                                          {
                                            "type": "Tag",
                                            "name": "h1",
                                            "selfClosing": false,
                                            "block": {
                                              "type": "Block",
                                              "nodes": [
                                                {
                                                  "type": "Code",
                                                  "val": "__('Welcome %s!', get('user.name'))",
                                                  "buffer": true,
                                                  "mustEscape": true,
                                                  "isInline": true,
                                                  "line": 5,
                                                  "filename": "pages/home.pug"
                                                }
                                              ],
                                              "line": 5,
                                              "filename": "pages/home.pug"
                                            },
                                            "attrs": [],
                                            "attributeBlocks": [],
                                            "isInline": false,
                                            "line": 5,
                                            "filename": "pages/home.pug"
                                          },
                                          {
                                            "type": "Code",
                                            "val": "Basket = get('user.basket')",
                                            "buffer": false,
                                            "mustEscape": false,
                                            "isInline": false,
                                            "line": 6,
                                            "filename": "pages/home.pug"
                                          },
                                          {
                                            "type": "Conditional",
                                            "test": "Basket",
                                            "consequent": {
                                              "type": "Block",
                                              "nodes": [
                                                {
                                                  "type": "Tag",
                                                  "name": "h3",
                                                  "selfClosing": false,
                                                  "block": {
                                                    "type": "Block",
                                                    "nodes": [
                                                      {
                                                        "type": "Text",
                                                        "val": "Basket:",
                                                        "line": 8,
                                                        "filename": "pages/home.pug"
                                                      }
                                                    ],
                                                    "line": 8,
                                                    "filename": "pages/home.pug"
                                                  },
                                                  "attrs": [],
                                                  "attributeBlocks": [],
                                                  "isInline": false,
                                                  "line": 8,
                                                  "filename": "pages/home.pug"
                                                },
                                                {
                                                  "type": "Tag",
                                                  "name": "ul",
                                                  "selfClosing": false,
                                                  "block": {
                                                    "type": "Block",
                                                    "nodes": [
                                                      {
                                                        "type": "Each",
                                                        "obj": "Basket",
                                                        "val": "Item",
                                                        "key": "i",
                                                        "block": {
                                                          "type": "Block",
                                                          "nodes": [
                                                            {
                                                              "type": "Tag",
                                                              "name": "li",
                                                              "selfClosing": false,
                                                              "block": {
                                                                "type": "Block",
                                                                "nodes": [
                                                                  {
                                                                    "type": "Tag",
                                                                    "name": "a",
                                                                    "selfClosing": false,
                                                                    "block": {
                                                                      "type": "Block",
                                                                      "nodes": [
                                                                        {
                                                                          "type": "Code",
                                                                          "val": "Item.Name",
                                                                          "buffer": true,
                                                                          "mustEscape": true,
                                                                          "isInline": true,
                                                                          "line": 11,
                                                                          "filename": "pages/home.pug"
                                                                        }
                                                                      ],
                                                                      "line": 11,
                                                                      "filename": "pages/home.pug"
                                                                    },
                                                                    "attrs": [
                                                                      {
                                                                        "name": "href",
                                                                        "val": "url('product.view', {Uid:Item.Uid})",
                                                                        "mustEscape": true
                                                                      }
                                                                    ],
                                                                    "attributeBlocks": [],
                                                                    "isInline": true,
                                                                    "line": 11,
                                                                    "filename": "pages/home.pug"
                                                                  }
                                                                ],
                                                                "line": 11,
                                                                "filename": "pages/home.pug"
                                                              },
                                                              "attrs": [
                                                                {
                                                                  "name": "class",
                                                                  "val": "'item'",
                                                                  "mustEscape": false
                                                                }
                                                              ],
                                                              "attributeBlocks": [],
                                                              "isInline": false,
                                                              "line": 11,
                                                              "filename": "pages/home.pug"
                                                            }
                                                          ],
                                                          "line": 11,
                                                          "filename": "pages/home.pug"
                                                        },
                                                        "line": 10,
                                                        "filename": "pages/home.pug"
                                                      }
                                                    ],
                                                    "line": 9,
                                                    "filename": "pages/home.pug"
                                                  },
                                                  "attrs": [],
                                                  "attributeBlocks": [],
                                                  "isInline": false,
                                                  "line": 9,
                                                  "filename": "pages/home.pug"
                                                }
                                              ],
                                              "line": 8,
                                              "filename": "pages/home.pug"
                                            },
                                            "alternate": null,
                                            "line": 7,
                                            "filename": "pages/home.pug"
                                          }
                                        ],
                                        "line": 4,
                                        "filename": "pages/home.pug"
                                      },
                                      "attrs": [
                                        {
                                          "name": "id",
                                          "val": "'hello-jade'",
                                          "mustEscape": false
                                        }
                                      ],
                                      "attributeBlocks": [],
                                      "isInline": false,
                                      "line": 4,
                                      "filename": "pages/home.pug"
                                    }
                                  ],
                                  "line": 9,
                                  "filename": "layouts/default.pug",
                                  "name": "content",
                                  "mode": "replace"
                                }
                              ],
                              "line": 8,
                              "filename": "layouts/default.pug"
                            },
                            "attrs": [
                              {
                                "name": "class",
                                "val": "'container'",
                                "mustEscape": false
                              }
                            ],
                            "attributeBlocks": [],
                            "isInline": false,
                            "line": 8,
                            "filename": "layouts/default.pug"
                          }
                        ],
                        "line": 8,
                        "filename": "layouts/default.pug"
                      },
                      "attrs": [
                        {
                          "name": "id",
                          "val": "'content'",
                          "mustEscape": false
                        },
                        {
                          "name": "class",
                          "val": "'page-content'",
                          "mustEscape": false
                        }
                      ],
                      "attributeBlocks": [],
                      "isInline": false,
                      "line": 8,
                      "filename": "layouts/default.pug"
                    }
                  ],
                  "line": 12,
                  "filename": "layouts/blank.pug",
                  "name": "body",
                  "mode": "replace"
                }
              ],
              "line": 11,
              "filename": "layouts/blank.pug"
            },
            "attrs": [
              {
                "name": "id",
                "val": "'page'",
                "mustEscape": false
              },
              {
                "name": "class",
                "val": "'page'",
                "mustEscape": false
              }
            ],
            "attributeBlocks": [],
            "isInline": false,
            "line": 11,
            "filename": "layouts/blank.pug"
          }
        ],
        "line": 2,
        "filename": "layouts/blank.pug"
      },
      "attrs": [
        {
          "name": "class",
          "val": "'no-js'",
          "mustEscape": false
        },
        {
          "name": "lang",
          "val": "get('site').Language",
          "mustEscape": true
        }
      ],
      "attributeBlocks": [],
      "isInline": false,
      "line": 2,
      "filename": "layouts/blank.pug"
    }
  ]
}
		`), "test.ast.json")
		Expect(p.render(testast, "", nil)).To(Equal(`<!DOCTYPE html>
<html class="no-js" lang="{{(get "site").Language}}">
  <body id="page" class="page">
	{{$Site := (get "site")}}<a name="top"></a>
	<section id="content" class="page-content">
	  <div class="container">
	    <div id="hello-jade">
		  <h1>{{(__ "Welcome %s!" (get "user.name"))}}</h1>
		  {{$Basket := (get "user.basket")}}
		  {{if $Basket}}
		  <h3>Basket:</h3>
		  <ul>
		    {{range $i, $Item := $Basket}}
		    <li class="item"><a href="{{(url "product.view" (__op__map "Uid" $Item.Uid))}}">{{$Item.Name}}</a></li>
		    {{end}}
		  </ul>
		  {{end}}
	    </div>
	  </div>
	</section>
  </body>
</html>`))
	})
})

func TestRender(t *testing.T) {
	//RegisterFailHandler(Fail)
	//RunSpecs(t, "Render Test Suite")
}
