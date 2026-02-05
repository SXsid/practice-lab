const setTimeoutPromisified = (delay) => {
  return new Promise((resolve, reject) => setTimeout(resolve, delay));
};

console.log(setTimeoutPromisified(200).then(() => console.log("Hello world!")));
