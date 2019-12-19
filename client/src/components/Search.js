import React from 'react';

import "./Search.less"

function Search({onSearch}) {

  const handleSubmit = (e) => {
    onSearch(e.target.elements.searchText)
  }

  return (
      <form className="search-form" onSubmit={handleSubmit}>
        <input name="searchText" type="text" />
        <input type="submit" />
      </form>
  );
}

export default Search;
