<h1 align="center">SSE Test 👍</h1>

[def]: ./frontend/sse.js

This project is a implementation of SSE or Server Side Events which represents a server-client unidirectional connection. It uses a Golang backend with a priority queue that manages the messages state to send logs through the clients/subcribers list. The frontend consumer is a simple html & js that connects to the golang server via EventSource class (see [code][def] for more understanding).
