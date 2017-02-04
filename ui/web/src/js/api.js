var $ = require('jquery');

class API {
  static formatPath() {
    var joined = Array.prototype.join.call(arguments, '/');
    return ('/' + joined).replace(/\/{2,}/g, '/');
  }

  constructor(basePath) {
    this.basePath = basePath;
  }

  getFeeds() {
    return $.get(API.formatPath(this.basePath, 'feeds'));
  }

  getFeed(id, cb) {
    return $.get(API.formatPath(this.basePath, 'feeds', id));
  }

  getItems(feedID) {
    return $.get(API.formatPath(this.basePath, 'feeds', feedID, 'items'));
  }

  getItem(feedID, itemID) {
    return $.get(API.formatPath(this.basePath, 'feeds', feedID, 'items', itemID));
  }
}

module.exports = API;
