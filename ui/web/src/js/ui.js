var React = require('react'),
  ReactDOM = require('react-dom'),
  App = require('./components/App'),
  log = require('./log');

var ui = {
  /*
   * Events generated by UI
   */
  eventHandlers: {
    // Select feed to see it's items
    onSelectFeed: [],
    onSubscribeToFeed: [],
    onSelectItem: []
  },

  /*
   * Attach a listener to one of UI events
   */
  registerHandler: function(key, fn) {
    if(ui.eventHandlers[key]) {
      ui.eventHandlers[key].push(fn);
    } else {
      log.err("Unknown handler", key);
    }
  },

  /*
   * Private
   *
   * UI events emitter
   */
  getHandler: function(key) {
    return function(data) {
      for(let h of ui.eventHandlers[key]) {
        h(data);
      }
    };
  },

  render: function(data) {
    ReactDOM.render(
      // jshint ignore:start
      <App feeds={data.feeds}
      onSelectFeed={this.getHandler('onSelectFeed')}
      onSubscribeToFeed={this.getHandler('onSubscribeToFeed')}
      onSelectItem={this.getHandler('onSelectItem')} />,
      // jshint ignore:end
      document.getElementById('root')
    );
  }
};

module.exports = ui;
