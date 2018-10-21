package main

func main() {
	a := App{}
	// You need to set your Username and Password here
	a.Initialize("tours_admin", "ladmdetouris", "restaurant")

	a.Run(":8080")
}
