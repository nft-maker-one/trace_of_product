const express = require('express');
const cors = require('cors');
const app = express();
const bodyParser = require('body-parser');
app.use(cors());
app.use(bodyParser.json());

let nodes = {
    "1304928106": "127.0.0.1:8082",
    "5020059120": "127.0.0.1:8085",
    "6554550297": "127.0.0.1:8084",
    "6711488489": "127.0.0.1:8083"
};

// Endpoint to get message by queryId
app.get('/message', (req, res) => {
    const queryId = req.query.id;
    if (queryId === '10000000003') {
        res.json({
            msg: "{\"eggplant_id\":10000000003,\"product_height\":1,\"product_hash\":[188,113,105,165,64,72,118,0,201,38,207,14,60,192,0,71,229,110,138,60,131,145,249,72,190,181,69,166,176,33,154,77],\"transport_height\":1,\"transport_hash\":[188,113,105,165,64,72,118,0,201,38,207,14,60,192,0,71,229,110,138,60,131,145,249,72,190,181,69,166,176,33,154,77],\"process_height\":0,\"process_hash\":[188,113,105,165,64,72,118,0,201,38,207,14,60,192,0,71,229,110,138,60,131,145,249,72,190,181,69,166,176,33,154,77],\"storage_height\":1,\"storage_hash\":[188,113,105,165,64,72,118,0,201,38,207,14,60,192,0,71,229,110,138,60,131,145,249,72,190,181,69,166,176,33,154,77],\"sell_height\":1,\"sell_hash\":[188,113,105,165,64,72,118,0,201,38,207,14,60,192,0,71,229,110,138,60,131,145,249,72,190,181,69,166,176,33,154,77]}"
        });
    } else {
        res.status(404).send('Not Found');
    }
});

// Endpoint to handle file uploads
app.post('/upload', (req, res) => {
    const data = req.body;

    // Print the received data
    console.log('Received data:', data);

    // TODO: Process the data, such as saving it to a database

    // Return success response
    res.json({ message: 'Data uploaded successfully', receivedData: data });
});

// New login endpoint
app.post('/login', (req, res) => {
    const { user_name, password } = req.body;

    // Hardcoded credentials for simplicity
    const validUserName = 'admin';
    const validPassword = 'agri_chain';

    // Check if the provided credentials match the valid ones
    if (user_name === validUserName && password === validPassword) {
        res.json({ success: true, message: 'Login successful' });
    } else {
        res.status(401).json({ success: false, message: 'Invalid username or password' });
    }
});

// Endpoint to get nodes
app.get('/nodes', (req, res) => {
    res.json(nodes);
});

// Function to generate a new node ID and address
function generateNewNode() {
    const newId = Date.now().toString();
    const newAddress = `127.0.0.1:${Math.floor(Math.random() * 10000) + 8000}`;
    nodes[newId] = newAddress;
}

// Add a new node every 2 seconds
setInterval(generateNewNode, 2000);

app.listen(8081, () => {
    console.log('Server is running on http://localhost:8081');
});

