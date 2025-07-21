import { useState } from 'react';
import './Login.css'

function Login() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })

  // Listen for typing
  const handleChange = ((event) =>{
    const{name, value} = event.target;
    setFormData((prevData) => ({
      ...prevData,
      [name] : value,
    }));
  }); 

  // Handle login form
  async function handleSubmit(event){
    event.preventDefault();
    try {
      const response = await fetch('http://localhost:5001/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          email: formData.email, 
          password: formData.password, 
        })
      });

      const data = await response.json();

       if (response.ok) {
        console.log('Login successful:', data);
        // Store token or user info
        localStorage.setItem('token', data.token);
        // Redirect or update UI
      } else {
        console.error('Login failed:', data.message);
      }
      setFormData({email:"", password:""});

    } catch(error){
      console.log(error);
    };
    
  }

  return (
    <div className="login-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <label htmlFor="email">Email</label>
        <input 
        type="email" 
        id="email" 
        name="email" 
        value={formData.email}
        onChange={handleChange}
        required/>

        <label htmlFor="password">Password</label>
        <input 
        type="password" 
        id="password" 
        name="password" 
        value={formData.password}
        onChange={handleChange}
        required/>

        <button type="submit">Login</button>
      </form>
    </div>
  );
}

export default Login
