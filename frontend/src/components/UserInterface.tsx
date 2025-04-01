import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CardComponent from './CardComponent';

interface User {
  id: number;
  name: string;
  statement: number;
}

interface UserInterfaceProps {
  backendName: string;
}

interface ValidationErrors {
  newStatement: string;
  updateId: string;
  updateStatement: string;
}

const UserInterface: React.FC<UserInterfaceProps> = ({ backendName }) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';
  const [users, setUsers] = useState<User[]>([]);
  const [newUser, setNewUser] = useState({ name: '', statement: '' });
  const [updateUser, setUpdateUser] = useState({ id: '', name: '', statement: '' });
  const [errors, setErrors] = useState<ValidationErrors>({
    newStatement: '',
    updateId: '',
    updateStatement: ''
  });

  const backgroundColors: { [key: string]: string } = {
    repair: 'bg-stone-900',
  };

  const buttonColors: { [key: string]: string } = {
    repair: 'bg-cyan-700 hover:bg-blue-600',
  };

  const bgColor = backgroundColors[backendName as keyof typeof backgroundColors] || 'bg-lime-500';
  const btnColor = buttonColors[backendName as keyof typeof buttonColors] || 'bg-gray-500 hover:bg-gray-600';

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/${backendName}/users`);
        setUsers(response.data.reverse());
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [backendName, apiUrl]);

  // validating input
  const validateNumberInput = (value: string): boolean => {
    return /^\d*$/.test(value);
  };

  const handleNewStatementChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    
    if (value === '' || validateNumberInput(value)) {
      setNewUser({ ...newUser, statement: value });
      setErrors({ ...errors, newStatement: '' });
    } else {
      setErrors({ ...errors, newStatement: 'Please use only numbers' });
    }
  };

  const handleUpdateIdChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    
    if (value === '' || validateNumberInput(value)) {
      setUpdateUser({ ...updateUser, id: value });
      setErrors({ ...errors, updateId: '' });
    } else {
      setErrors({ ...errors, updateId: 'Please use only numbers' });
    }
  };

  const handleUpdateStatementChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    
    if (value === '' || validateNumberInput(value)) {
      setUpdateUser({ ...updateUser, statement: value });
      setErrors({ ...errors, updateStatement: '' });
    } else {
      setErrors({ ...errors, updateStatement: 'Please use only numbers' });
    }
  };

  const createUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${apiUrl}/api/${backendName}/users`, { ...newUser, statement: Number(newUser.statement) });
      setUsers([response.data, ...users]);
      setNewUser({ name: '', statement: '' });
    } catch (error) {
      console.error('Error creating user:', error);
    }
  };

  const handleUpdateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const currentName = updateUser.name || users.find(user => user.id === parseInt(updateUser.id))?.name || '';
      await axios.put(`${apiUrl}/api/${backendName}/users/${updateUser.id}`, { name: currentName, statement: Number(updateUser.statement) });
      setUpdateUser({ id: '', name: '', statement: '' });
      setUsers(
        users.map((user) => {
          if (user.id === parseInt(updateUser.id)) {
            return { ...user, name: currentName, statement: Number(updateUser.statement) };
          }
          return user;
        })
      );
    } catch (error) {
      console.error('Error updating user:', error);
    }
  };

  const deleteUser = async (userId: number) => {
    try {
      await axios.delete(`${apiUrl}/api/${backendName}/users/${userId}`);
      setUsers(users.filter((user) => user.id !== userId));
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  }

  return (
    <div className={`user-interface ${bgColor} ${backendName} w-full max-w-md p-4 my-4 rounded shadow`}>
      <img src={`/${backendName}logo.svg`} alt={`${backendName} Logo`} className="w-20 h-20 mb-6 mx-auto" />
      <h2 className="text-xl font-bold text-center text-white mb-6">{`${backendName.charAt(0).toUpperCase() + backendName.slice(1)} Fund, bl.4E, ul. Dragoman, Sofia`}</h2>

      {/* Create user */}
      <form onSubmit={createUser} className="mb-6 p-4 bg-blue-100 rounded shadow">
        <input
          placeholder="Name"
          value={newUser.name}
          onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
          className="mb-2 w-full p-2 border border-gray-300 rounded"
        />
        <div className="mb-2 relative">
          <input
            placeholder="Statement"
            value={newUser.statement}
            onChange={handleNewStatementChange}
            className={`w-full p-2 border rounded ${errors.newStatement ? 'border-red-500' : 'border-gray-300'}`}
          />
          {errors.newStatement && (
            <p className="text-red-500 text-xs mt-1">{errors.newStatement}</p>
          )}
        </div>
        <button type="submit" className="w-full p-2 text-white bg-blue-500 rounded hover:bg-blue-600">
          Add User
        </button>
      </form>

      {/* Update user */}
      <form onSubmit={handleUpdateUser} className="mb-6 p-4 bg-blue-100 rounded shadow">
        <div className="mb-2 relative">
          <input
            placeholder="User Id"
            value={updateUser.id}
            onChange={handleUpdateIdChange}
            className={`w-full p-2 border rounded ${errors.updateId ? 'border-red-500' : 'border-gray-300'}`}
          />
          {errors.updateId && (
            <p className="text-red-500 text-xs mt-1">{errors.updateId}</p>
          )}
        </div>
        <input
          placeholder="New Name"
          value={updateUser.name}
          onChange={(e) => setUpdateUser({ ...updateUser, name: e.target.value })}
          className="mb-2 w-full p-2 border border-gray-300 rounded"
        />
        <div className="mb-2 relative">
          <input
            placeholder="New statement"
            value={updateUser.statement}
            onChange={handleUpdateStatementChange}
            className={`w-full p-2 border rounded ${errors.updateStatement ? 'border-red-500' : 'border-gray-300'}`}
          />
          {errors.updateStatement && (
            <p className="text-red-500 text-xs mt-1">{errors.updateStatement}</p>
          )}
        </div>
        <button type="submit" className="w-full p-2 text-white bg-green-500 rounded hover:bg-green-600">
          Update User
        </button>
      </form>

      {/* display users */}
      <div className="space-y-4">
        {users.map((user) => (
          <div key={user.id} className="flex items-center justify-between bg-white p-4 rounded-lg shadow">
            <CardComponent card={user} />
            <button onClick={() => deleteUser(user.id)} className={`${btnColor} text-white py-2 px-4 rounded`}>
              Delete User
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default UserInterface;