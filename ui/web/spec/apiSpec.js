var API = require('../src/js/api');

describe('API', function() {
  describe('.formatPath', function() {
    var fn = API.formatPath;

    it('handles slashes correctly', function() {
      expect(fn('//', 'a', '///', 'b', 'c')).toBe('/a/b/c');
    });
  });

  describe('ajax methods', function() {
    beforeEach(function() {
      jasmine.Ajax.install();
    });

    afterEach(function() {
      jasmine.Ajax.uninstall();
    });

    describe('getFeeds', function() {
      it('gets the feeds from expected url', function() {
      });
    });
  });
});
