import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Header } from './components/header';
import { WebSocketProvider } from './context/websocket';
import { Dashboard } from './pages/dashboard/dashboard';

function App() {
  return (
    <BrowserRouter>
      <WebSocketProvider>
        <main className="relative w-full h-screen flex flex-col overflow-hidden">
          <Header />
          <Routes>
            <Route path="/" element={<Dashboard />} />
          </Routes>
        </main>
      </WebSocketProvider>
    </BrowserRouter>
  );
}

export default App;
