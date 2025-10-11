export const ThrothelJS = () => {
  resize();
};

const throthel = (handler, delay) => {
  let timer = 0;
  return (...args) => {
    if (timer) return;
    timer = setTimeout(() => {
      handler(...args);
      timer = 0;
    }, delay);
  };
};
const rezieHandler = () => {
  //dummy sving
  console.log("rezing happeing");
};
const resize = () => {
  const box = document.querySelector(".box");
  const throthelRezier = throthel(rezieHandler, 2000);
  const observer = new ResizeObserver(() => throthelRezier());

  observer.observe(box);
};
