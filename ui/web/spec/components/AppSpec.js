import React from "react";
import ReactDOM from 'react-dom';
import App from '../../src/js/components/App';

describe('App', function() {
  var feeds = [
    {ID: 1, Title: "The Title", Link: "http://example.com"}
  ],
    onSelectFeed = function() {},
    root = document.createElement('div');

  function render() {
    ReactDOM.render(
      // jshint ignore:start
      <App feeds={feeds} onSelectFeed={onSelectFeed} />,
      // jshint ignore:end
      root
    );
  }

  beforeEach(function() {
    document.querySelector('body').appendChild(root);
  });

  afterEach(function() {
    document.querySelector('body').removeChild(root);
  });

  it('is rendered without troubles', function() {
    expect(root.children.length).toEqual(0);
    render();
    expect(root.children.length).toEqual(1);
  });
});
