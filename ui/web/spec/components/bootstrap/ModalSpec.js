require('../../../src/index');
import Modal from "../../../src/js/components/bootstrap/Modal";

describe('Modal', function() {
	var root, onHide, title, body;

	title = "The Title";
	body = "The Body";
	root = document.createElement('div');
	onHide = function() {};

	function render() {
		ReactDOM.render(
			// jshint ignore:start
			<Modal show={true} onHide={onHide} title={title} body={body} />,
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
		expect(root.children.length).toBe(0);
		render();
		expect(root.children.length).toBe(1);
	});

	it('calls onHide when modal was hidden because user clicked "close" button', function() {
		onHide = jasmine.createSpy('onHide');
		render();
		root.querySelector('[data-dismiss="modal"]').click();
		expect(onHide).toHaveBeenCalled();
	});

	it('renders title and body', function() {
		render();
		expect(root.querySelector('.modal-title').innerText).toEqual(title);
		expect(root.querySelector('.modal-body').innerText).toEqual(body);
	});
});
