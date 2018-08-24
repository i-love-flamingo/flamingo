/*
This package provides helper to use and deal with (web)forms in your "interfaces" layer (e.g. in your controllers)
(Use this package only in our "interface" layer (e.g. in your controllers) of your own module.)

Usage:

- Add your Data Representation of your form to your package e.g. in the folder ("/interfaces/controller/form")

- To process your form use "SimpleProcessFormRequest" or "ProcessFormRequest"

- In case you want to use "ProcessFormRequest", you need to write an implementation of the interface "domain.FormService"

- Optional your implementation can also implement the interface "domain.GetDefaultFormData", to be able to prepopulate your form data

Details

Check the documentation of the application and form package
*/
package form
