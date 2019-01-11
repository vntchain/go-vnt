package main

import (
	"bytes"
)

var buffer = bytes.NewBuffer([]byte(`
[{
     "name": "xxxx",
     "tables": [{
          "name": "xxxx",
          "type": "",
          "tables": [{
                    "name": "key",
                    "type": "int32",
                    "tables": []
               },
               {
                    "name": "value",
                    "type": "",
                    "tables": [{
                              "name": "key",
                              "type": "int32",
                              "tables": []
                         },
                         {
                              "name": "value",
                              "type": "",
                              "tables": [{
                                        "name": "map1",
                                        "type": "",
                                        "tables": [{
                                                  "name": "key",
                                                  "type": "int32",
                                                  "tables": []
                                             },
                                             {
                                                  "name": "value",
                                                  "type": "int32",
                                                  "tables": []
                                             }
                                        ]
                                   },
                                   {
                                        "name": "array1",
                                        "type": "",
                                        "tables": [{
                                                  "name": "index",
                                                  "type": "uint64",
                                                  "tables": []
                                             },
                                             {
                                                  "name": "value",
                                                  "type": "",
                                                  "tables": [{
                                                            "name": "value",
                                                            "type": "int32",
                                                            "tables": []
                                                       },
                                                       {
                                                            "name": "key",
                                                            "type": "int32",
                                                            "tables": []
                                                       }
                                                  ]
                                             },
                                             {
                                                  "name": "length",
                                                  "type": "uint64",
                                                  "tables": []
                                             }
                                        ]
                                   },
                                   {
                                        "name": "name",
                                        "type": "string",
                                        "tables": []
                                   },
                                   {
                                        "name": "address",
                                        "type": "string",
                                        "tables": []
                                   },
                                   {
                                        "name": "phone",
                                        "type": "string",
                                        "tables": []
                                   }
                              ]
                         }
                    ]
               }
          ]
     }],
     "type": "key"
}]


`))
