package models

// Model relationships are defined here to avoid circular dependencies
// These types are used in the model structs to represent relationships

// SubscriptionRelations represents the relationships for a Subscription
type SubscriptionRelations struct {
	Orders []*Order
}

// OrderRelations represents the relationships for an Order
type OrderRelations struct {
	Subscription    *Subscription
	Item            *Item
	BillingSchedule *BillingSchedule
}

// InvoiceRelations represents the relationships for an Invoice
type InvoiceRelations struct {
	Subscription *Subscription
	LineItems    []*InvoiceLineItem
}

// InvoiceLineItemRelations represents the relationships for an InvoiceLineItem
type InvoiceLineItemRelations struct {
	Invoice *Invoice
	Order   *Order
}
