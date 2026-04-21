import { Routes, Route } from 'react-router-dom';
import Sidebar from './components/Sidebar/Sidebar';
import { ThemeProvider } from './context/ThemeContext';
import { InvoiceProvider } from './context/InvoiceContext';
import './App.css';
import InvoiceListPage from './pages/InvoiceListPage/InvoiceListPage';
import InvoiceDetailPage from './pages/InvoiceDetailPage/InvoiceDetailPage';

function App() {
  return (
    <ThemeProvider>
      <InvoiceProvider>
        <div className="app-container">
          <Sidebar />
          <main className="main-content">
            <Routes>
              <Route path="/" element={<InvoiceListPage />} />
              <Route path="/invoice/:id" element={<InvoiceDetailPage />} />
            </Routes>
          </main>
        </div>
      </InvoiceProvider>
    </ThemeProvider>
  );
}

export default App;
