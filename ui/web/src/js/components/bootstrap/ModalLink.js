import React from "react";
import Modal from "./Modal";

class ModalLink extends React.Component {
	constructor(props) {
		super(props);
		this.show = this.show.bind(this);
		this.modalWasHidden = this.modalWasHidden.bind(this);
		this.state = { showModal: false };
	}

	show(e) {
		e.preventDefault();
		e.stopPropagation();
		this.setState({showModal: true});
	}

	modalWasHidden() {
		this.setState({showModal: false});
	}

	render() {
		var modal = null;

		if(this.state.showModal) {
			// jshint ignore:start
			modal = <Modal show={true} onHide={this.modalWasHidden} title={this.props.title}
										 body={this.props.body} />
			// jshint ignore:end
		}

		return (
			// jshint ignore:start
			<div>
				<a href="#" onClick={this.show}>{this.props.linkText}</a>
				{modal}
			</div>
			// jshint ignore:end
		);
	}
}

ModalLink.propTypes = {
	linkText: React.PropTypes.string
};

module.exports = ModalLink;
