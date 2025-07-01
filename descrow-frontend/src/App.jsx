
import './App.css'

function App() {

  return (
    <div className="login-container">
      <h2>Login</h2>
      <form>
        <label htmlFor="email">Email</label>
        <input type="email" id="email" name="email" required/>

        <label htmlFor="password">Password</label>
        <input type="password" id="password" name="password" required/>

        <button type="submit">Login</button>
      </form>
    </div>
  );
}

export default App
