import { useEffect, useMemo, useState } from "react";
import axios from "axios";
import "./App.css";

function App() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    const fetchCurrencies = async () => {
      try {
        const response = await axios.get("http://localhost:8080/api/currencies/today");
        setData(response.data);
      } catch (err) {
        setError("Veri alınamadı.");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCurrencies();
  }, []);

  const highlightedCurrencies = useMemo(() => {
    if (!data) return [];

    const priorityCodes = ["USD", "EUR", "GBP"];
    return priorityCodes
        .map((code) => data.currencies.find((currency) => currency.code === code))
        .filter(Boolean);
  }, [data]);

  const filteredCurrencies = useMemo(() => {
    if (!data) return [];

    return data.currencies.filter((currency) => {
      const value =
          `${currency.code} ${currency.currencyName} ${currency.currencyNameTr}`.toLowerCase();

      return value.includes(searchTerm.toLowerCase());
    });
  }, [data, searchTerm]);

  if (loading) {
    return (
        <div className="container">
          <h1>Yükleniyor...</h1>
        </div>
    );
  }

  if (error) {
    return (
        <div className="container">
          <h1>{error}</h1>
        </div>
    );
  }

  return (
      <div className="container">
        <header className="page-header">
          <div>
            <h1>TCMB Currency Dashboard</h1>
            <p className="subtitle">Güncel döviz kurları</p>
          </div>
          <div className="meta">
            <p><strong>Tarih:</strong> {new Date(data.date).toLocaleDateString("tr-TR")}</p>
            <p><strong>Bülten No:</strong> {data.dayNo}</p>
          </div>
        </header>

        <section className="cards">
          {highlightedCurrencies.map((currency) => (
              <div className="card" key={currency.code}>
                <h2>{currency.code}</h2>
                <p>{currency.currencyNameTr}</p>
                <div className="card-values">
                  <span>Alış: {currency.forexBuying}</span>
                  <span>Satış: {currency.forexSelling}</span>
                </div>
              </div>
          ))}
        </section>

        <section className="toolbar">
          <input
              type="text"
              placeholder="Kod veya para birimi ara..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
          />
        </section>

        <section className="table-wrapper">
          <table>
            <thead>
            <tr>
              <th>Kod</th>
              <th>Ad</th>
              <th>Türkçe Ad</th>
              <th>Birim</th>
              <th>Forex Alış</th>
              <th>Forex Satış</th>
              <th>Efektif Alış</th>
              <th>Efektif Satış</th>
            </tr>
            </thead>
            <tbody>
            {filteredCurrencies.map((currency) => (
                <tr key={currency.code}>
                  <td>{currency.code}</td>
                  <td>{currency.currencyName}</td>
                  <td>{currency.currencyNameTr}</td>
                  <td>{currency.unit}</td>
                  <td>{currency.forexBuying}</td>
                  <td>{currency.forexSelling}</td>
                  <td>{currency.banknoteBuying}</td>
                  <td>{currency.banknoteSelling}</td>
                </tr>
            ))}
            </tbody>
          </table>
        </section>
      </div>
  );
}

export default App;