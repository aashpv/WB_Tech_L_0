async function fetchOrderDetails() {
    const id = document.getElementById('id').value;
    const response = await fetch(`http://localhost:8090/api/${id}`);

    if (response.ok) {
        const data = await response.json();

        let orderDetailsHtml = `<h2>Order ${id} Details</h2>`;
        orderDetailsHtml += "<pre>" + JSON.stringify(data, null, 2) + "</pre>";

        document.getElementById('orderDetails').innerHTML = orderDetailsHtml;
    } else if (response.status === 404) {
        document.getElementById('orderDetails').innerHTML = "<p>Order not found</p>";
    } else {
        document.getElementById('orderDetails').innerHTML = "<p>Something went wrong</p>";
    }
}