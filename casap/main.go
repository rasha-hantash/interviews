// To execute Go code, please declare a func main() in a package "main"

/*

 Design a shopping cart service for a grocery app. You can define all the apis, entities and events you need to.


 A user should be able to add items from a menu to the shopping cart, checkout the shopping cart/mark the shopping cart as ready. The user can add items up until the shopper has started shopping for the items. When an item that is requested is unavailable, allow the shopper to suggest a replacement and provide the user a chance to either approve, reject or suggest another replacement. If no replacement is available then allow the shopper to refund the item/mark the item as not found. Then when the shopping is done, shoper should say shopping is done, the items should not be editable by anyone.

 Timeline: 
 
 User checks out shopping cart (virtual cart, marks as ready to shop) then we dispatch shopper to grocery store. Shopper reaches store and says started shapping, marking items as found, not found or replaced. Then user approves/rejects/suggests replacements. Shopping is done.
 */

 package main

 import "fmt"
 
 type ReplacementResponseStatus int
 
 const (
	 ReplacementResponseStatusAccepted ReplacementResponseStatus = iota
	 ReplacementResponseStatusRejected
	 ReplacementResponseStatusUserSuggestedAlternative
 )
 
 // Item 1 has been found
 // Item 2 has been replaces is waiting on user
 // Item 3 - still needs to worked on
 
 /*
	 Items -> id, name, alternative, shopping_cart_id , status 
	 ShoppingCarts -> id, status, items
 
	 select status from items where shopping_cart_id == $id; 
 */
 
 type Item struct {
	 ID string 
	 ItemName string 
	 Alternative string
	 Status string `db:"status"`
 }
 
 type ShoppingCart struct {
	 ID string 
	 status string // is_ready , etc. 
	 Items []Item 
 }
 
 
 type ShoppingCartService interface {
	 AddItem(item Item) (string, error)
   Checkout(shoppingCartId string) (string, error)
   UpdateStatus(status string) error
	 ProposeReplacement(item Item) (ReplacementResponseStatus)// 
	 AcceptReplacement()
	 RejectReplacement() // process refund
	 SuggestAlternative() 
	 
 }
 
 
 
 func main() {
   for i := 0; i < 5; i++ {
	 fmt.Println("Hello, World!")
   }
 }
 