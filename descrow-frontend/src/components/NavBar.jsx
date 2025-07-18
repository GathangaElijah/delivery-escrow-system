import SearchBar from "./SearchBar";
import "./NavBar.css";

function NavBar({ onSearch }){
  return (
    <nav className="navbar">
      <div className="logo">MyApp</div>
      
      <div className="navbar-right">
        <SearchBar onSearch={onSearch} />
        {/* Future components like profile, icons, etc will go here */}
      </div>
    </nav>
  );
}
export default NavBar;