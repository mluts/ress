module.exports = {
  log: function() {
    window.console.log.apply(window.console, arguments);
  },

  err: function() {
    window.console.error.apply(window.console, arguments);
  }
};
