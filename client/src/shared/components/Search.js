import React, { useRef, useEffect, useState } from "react";
import { ReactComponent as SearchIcon } from "@shared/icons/search.svg";
import { ReactComponent as CaretDownIcon } from "@shared/icons/caret-down.svg";
import { useHideOnOutsideClick } from "@shared/hooks";
import "./Search.scss";

function Search({ onSearch, sources }) {
  const inputRef = useRef();
  const sourceDropdownRef = useRef();
  const [shouldShowSourceDropdown, setShouldShowSourceDropdown] = useState(false);

  const handleSubmit = e => {
    e.preventDefault();
    onSearch({
      searchText: e.target.elements.searchText.value,
      sources: sources.filter(s => e.target.elements["source-" + s.id].checked).map(s => s.id),
    });
  };

  const showSourceSelectDropdown = () => setShouldShowSourceDropdown(true);
  const hideSourceSelectDropdown = () => setShouldShowSourceDropdown(false);

  useEffect(() => {
    inputRef.current.focus();
  }, []);

  useHideOnOutsideClick(sourceDropdownRef, hideSourceSelectDropdown);

  return (
    <form className="search-form" onSubmit={handleSubmit}>
      <button className="show-sources" type="button" onClick={showSourceSelectDropdown}>
        <CaretDownIcon />
      </button>
      <input ref={inputRef} name="searchText" type="text" autoComplete="off" />
      <button type="submit">
        <SearchIcon />
      </button>

      <ul className={`${shouldShowSourceDropdown ? "" : "hide"} source-select-dropdown`} ref={sourceDropdownRef}>
        {sources.map(source => (
          <li key={source.id}>
            <label>
              <input type="checkbox" name={"source-" + source.id} value={source.id} defaultChecked></input>
              <span>{source.name}</span>
            </label>
          </li>
        ))}
      </ul>
    </form>
  );
}

export default Search;
