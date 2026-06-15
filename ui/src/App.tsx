import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Header } from "./components/header";
import { Dashboard } from "./pages/dashboard/dashboard";
import AgentDetails from "./pages/agentDetails/agentDetails";
import { useEffect } from "react";
import { connectWS } from "./services/websocket/agentSocket";

function App() {
  useEffect(() => {
    connectWS();
  }, []);

  return (
    <BrowserRouter>
      <main className="relative w-full h-screen flex flex-col overflow-hidden">
        <Header />
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/agents/:guid" element={<AgentDetails />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

export default App;
