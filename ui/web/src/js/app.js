var API = require('./api'),
    React = require('react'),
    ReactDOM = require('react-dom'),
    Sidebar = require('./sidebar'),
    log = require('./log');

function renderSidebar(feeds, onFeedSelected) {
  ReactDOM.render(
    // jshint ignore:start
    <Sidebar onFeedSelected={onFeedSelected} feeds={feeds} />,
    // jshint ignore:end
    document.getElementById('root')
  );
}

module.exports = {
  main: function() {
    var api = new API('/api');

    api.getFeeds().then(feeds => {
      if(_.isArray(feeds)) {
        renderSidebar(feeds, f => { console.log(f); });
      } else {
        log.err("Bad feeds response:", feeds);
      }
    });
  }
};
