import { useState, type ChangeEvent, type ChangeEventHandler } from "react";
import "./App.css";

interface todo {
  title: string;
  desc: string;
  isComplete: boolean;
}
function App() {
  const [todo, setTodo] = useState<Array<todo>>([]);
  const [Filtertodo, setFilteredTodo] = useState<Array<todo>>([]);
  const [currentTodo, setCurrentTodo] = useState<todo>({
    title: "",
    desc: "",
    isComplete: false,
  });
  const handeInput = (e: ChangeEvent<HTMLInputElement>) => {
    setCurrentTodo((value) => {
      return { ...value, [e.target.name]: e.target.value };
    });
  };
  const handleAdd = () => {
    setTodo((prevTodo) => {
      return [...prevTodo, currentTodo];
    });
  };
  const handleChange = (e: ChangeEvent<HTMLSelectElement>) => {
    const filterValue = e.target.value;
    setFilteredTodo(() => {
      if (filterValue === "All") return todo;
      return todo.filter((value) => String(value.isComplete) == e.target.value);
    });
  };
  return (
    <>
      <h1 className={"heading"}>TODO App</h1>
      {currentTodo.title}
      {currentTodo.desc}
      <div className={"main"}>
        <div className={"inputbox"}>
          <input
            name="title"
            value={currentTodo.title}
            type="text"
            placeholder="title"
            onChange={(e) => handeInput(e)}
          />

          <textarea
            name="desc"
            value={currentTodo.desc}
            onChange={(e) => handeInput(e)}
            placeholder="description"
          />
          <button onClick={handleAdd}>ADD</button>
        </div>
      </div>
      <select name="filter" id="" onChange={handleChange}>
        <option value="All">All</option>
        <option value={"true"}>completed</option>
        <option value="false">not-completed</option>
      </select>
      <div>
        {Filtertodo?.map((value) => {
          return (
            <div>
              <h2>{value.title}</h2>
              <h3>{value.desc}</h3>
            </div>
          );
        })}
      </div>
    </>
  );
}

export default App;
