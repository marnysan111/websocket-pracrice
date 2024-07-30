import { useState, useEffect, useRef } from 'react';

function App() {
  const [response, setResponse] = useState([]);
  const [chat, setChat] = useState([]);
  const [connected, setConnected] = useState(false);
  const [text, setText] = useState("");
  const [roomID, setRoomID] = useState("");
  const ws = useRef(null);

  useEffect(() => {
    fetch("http://localhost:8080/api")
      .then(res => res.json())
      .then(data => {
        setResponse(data.message);
      });
  }, []);

  const handleTextInput = (e) => {
    setText(e.target.value);
  };

  const handleButton = () => {
    if (ws.current && connected) {
      ws.current.send(text);
      setText("");
    }
  };

  const handleInput = (e) => {
    setRoomID(e.target.value);
  };

  const connRoom = () => {
    if (!connected) {
      ws.current = new WebSocket(`ws://localhost:8080/ws/${roomID}`);
      
      ws.current.onopen = () => {
        console.log('WebSocket Client Connected');
        setConnected(true);
      };

      ws.current.onmessage = (message) => {
        console.log('Received message:', message.data);
        setChat(prevChat => [...prevChat, `Server: ${message.data}`]);
      };

      ws.current.onerror = (error) => {
        console.error('WebSocket Error: ', error);
      };

      ws.current.onclose = (event) => {
        console.log('WebSocket Client Disconnected: ', event.reason);
        setConnected(false);
        setChat([])
      };
    }
  };

  const disconnectWebSocket = () => {
    if (connected && ws.current) {
      ws.current.close();
    }
  };

  return (
    <div className="flex flex-col items-center p-4 bg-gray-100 min-h-screen">
      <div>{roomID} {connected ? "Connected" : "Disconnected"}</div>
      <input 
        onChange={handleInput}
        placeholder='Room ID'
        value={roomID}
        className="p-2 border border-gray-300 rounded mb-2"
      />
      <button onClick={connRoom} className="p-2 bg-blue-500 text-white rounded mb-2">
        接続
      </button>
      <button onClick={disconnectWebSocket} className="p-2 bg-red-500 text-white rounded mb-4">
        切断
      </button>
      <div className="w-full max-w-lg border border-gray-300 m-4 p-4 rounded-2xl bg-white shadow-lg">
        <textarea
          rows="4"
          value={text}
          className="w-full p-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          onChange={handleTextInput}
        ></textarea>
      </div>
      <button 
        onClick={handleButton} 
        className='text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 mb-4'>
        送信
      </button>
      <div className="w-full max-w-lg mt-4 p-4 bg-white border border-gray-300 rounded-2xl shadow-lg">
        RoomID: {roomID}
        {chat.map((message, index) => (
          <div key={index} className="p-2 border-b border-gray-200">{message}</div>
        ))}
      </div>
    </div>
  );
}

export default App;
