import { useEffect } from "react";

/**
 * Hook that alerts clicks outside of the passed ref
 */
export default function useHideOnOutsideClick(ref, onOutsideClick) {
  /**
   * Alert if clicked on outside of element
   */
  const handleClickOutside = event => {
    if (ref.current && !ref.current.contains(event.target)) {
      onOutsideClick();
    }
  };

  useEffect(() => {
    // Bind the event listener
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      // Unbind the event listener on clean up
      document.removeEventListener("mousedown", handleClickOutside);
    };
  });
}
