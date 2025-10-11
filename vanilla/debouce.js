export const debounceJS = () => {
  typing();
};

function typing() {
  const inputBox = document.querySelector(".input");
  const debouncedHandler = Debounce(typingHandler, 200);
  inputBox.addEventListener("input", debouncedHandler);
}

function Debounce(handler, delay) {
  let timer = 0;
  return (e) => {
    if (timer) {
      clearTimeout(timer);
    }

    timer = setTimeout(() => handler(e), delay);
  };
}

function typingHandler(e) {
  //dummy api request
  console.log(e.target.value);
}
