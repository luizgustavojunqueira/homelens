import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Header } from "./components/header";
import { Dashboard } from "./pages/dashboard/dashboard";
import AgentDetails from "./pages/agentDetails/agentDetails";
import { useEffect } from "react";
import { connectWS } from "./services/websocket/agentSocket";
import { ToastContainer } from "react-toastify";
import Settings from "./pages/settings/Settings";

function App() {
  useEffect(() => {
    connectWS();
  }, []);

  return (
    <>
      <BrowserRouter>
        <main className="relative w-full h-screen flex flex-col overflow-hidden">
          <Header />
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/agents/:guid" element={<AgentDetails />} />
            <Route path="/settings/" element={<Settings />} />
          </Routes>
        </main>
      </BrowserRouter>
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={true}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="dark"
      />
    </>
  );
}

export default App;
