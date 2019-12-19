import React, { useState } from "react";
import Search from "./components/Search";

import "./App.less";

function App() {
  const [searchText, setSearchText] = useState("");
  return (
    <div className="app">
      <Search onSearch={setSearchText} />
      <section>{searchText}</section>
    </div>
  );
}

export default App;
