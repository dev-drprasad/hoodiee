import React from "react";

import "./Search.scss";

function Search({ onSearch }) {
  const handleSubmit = e => {
    e.preventDefault();
    onSearch(e.target.elements.searchText.value);
  };

  return (
    <form className="search-form" onSubmit={handleSubmit}>
      <input name="searchText" type="text" />
      <input type="submit" />
    </form>
  );
}

export default Search;
