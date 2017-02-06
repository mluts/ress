import ContextMenu from "../../src/js/components/ContextMenu";

describe('ContextMenu', function() {
  var callback, root;

  const options = {
    a: "Option A",
    b: "Option B",
    c: "Option C"
  };

  const top = 100, left = 100;

  function render() {
    ReactDOM.render(
      // jshint ignore:start
      <ContextMenu top={top} left={left} options={options}
      onSelectOption={callback} show={true} />,
      // jshint ignore:end
      root
    );
  }

  beforeEach(function() {
    root = document.createElement('div');
    document.querySelector('body').appendChild(root);
    callback = jasmine.createSpy("callback");

    render();
  });

  it("renders div with options", function() {
    var list = root.querySelectorAll('div > ul > li');
    expect(list.length).toEqual(3);
    expect(list[0].innerText).toBe(options.a);
    expect(list[1].innerText).toBe(options.b);
    expect(list[2].innerText).toBe(options.c);
  });

  it("calls the callback", function() {
    root.querySelector('li').click();
    expect(callback).toHaveBeenCalledWith(
      Object.keys(options)[0]
    );
  });

  it('hides after first click', function() {
    root.querySelector('li').click();
    expect(root.querySelectorAll('li').length).toEqual(0);
  });

  it('works if was hidden once', function() {
    root.querySelector('li').click();
    render();
    expect(root.querySelectorAll('li').length).toBe(3);
  });
});
