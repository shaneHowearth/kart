import './App.css';
import { useState, useEffect } from 'react';

const API_URL = 'http://localhost:8080/api/product';
const API_ENDPOINT = 'http://localhost:8080/api/order';

function App() {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [cart, setCart] = useState([]);

    // fetch all the products
    useEffect(() => {
        fetch(API_URL)
            .then(response =>  response.json())
            .then(data => {
                setProducts(data);
                setLoading(false);
            })
            .catch(error => {
                console.error('Error fetching products:', error);
                setLoading(false);
            });
    }, []);

    // Create a function that adds a product to the cart.
    const addToCart = (product) => {
        setCart(prevCart =>{
            const existing = prevCart.find(item => item.id === product.id);
            if (existing) {
                return prevCart.map(item =>
                    item.id === product.id
                    ? { ...item, quantity: item.quantity +1 }
                    : item
                );
            }
            return [...prevCart, { ...product, quantity: 1 }];
        });
    }

    // Create a function that posts an order to the server.
    const makeOrder = async (cart) => {
        try {
            // Transform cart items to match API format
            const orderItems = cart.map(item => ({
                productId: item.id,
                quantity: item.quantity
            }));
            const response = await fetch(API_ENDPOINT, {
                method: 'POST', // Specify the HTTP method as POST
                headers: {
                    // Tell the server we are sending JSON data
                    'Content-Type': 'application/json',
                    // Optional: API Key or other authorization headers would go here
                },
                // Convert the JavaScript object (formData) into a JSON string
                body: JSON.stringify({
                    items: orderItems,
                }),
            });

            // Check if the request was successful (status code 200-299)
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Parse the JSON response from the server
            const data = await response.json();

            // Success handling
            // setMessage(`Success! Data posted with ID: ${data.id}. (View console for full response)`);
            // setIsSuccess(true);
            console.log('Server Response:', data);
        } catch (error) {
            // Error handling
            // setMessage(`Failed to post data: ${error.message}.`);
            // setIsSuccess(false);
            console.error('POST Error:', error);
        } finally {
            setLoading(false);
        }
    }

    if (loading) {
        return <div>Loading products...</div>;
    }

    return (
        <div className="App">
            <h1>Shane's Awesome Shopping Kart</h1>

            <div className="cart-section">
                <h2>Your Cart</h2>
                {cart.length === 0 ? (
                    <p>Your cart is empty</p>
                ) : (
                    <div className="cart-items">
                        {cart.map(item => (
                            <div key={item.id} className="cart-item">
                                <span>{item.name}</span>
                                <span>x{item.quantity}</span>
                                <span>{item.price}</span>
                            </div>
                        ))}
                        <button onClick={() => makeOrder(cart)}>Make Order</button>
                    </div>
                )}
            </div>
            <p>Products</p>
            <div className="product-list">
                {products.map(product => (
                    <div key={product.id} className="product-card">
                        <h3>{product.name}</h3>
                        <p className="category">{product.category}</p>
                        <p className="price">{product.price}</p>
                        <button onClick={() => addToCart(product)}>Add to Cart</button>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default App;
