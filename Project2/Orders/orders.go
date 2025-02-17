package Orders

type Order struct {
	/*
			    "user_id": 123,
		    "product_id": 456,
		    "quantity": 2,
		    "total_price": 500.00

	*/
	User_Id     int     `json:"user_id"`
	Product_Id  int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	Total_Price float64 `json:"total_price"`
}

var OrderList []Order
