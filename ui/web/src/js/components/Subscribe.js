import React from "react";
import Modal from './bootstrap/Modal';
import _ from "lodash";

import Input from "./bootstrap/Input";
import Form from "./bootstrap/Form";
import Submit from "./bootstrap/Submit";

class Subscribe extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      Link: "",
      showModal: false
    };

    _.bindAll(this, 'handleChange', 'handleSubmit',
      'handleModalHide', 'handleClick');
  }

  handleChange(e) {
    this.setState({
      [e.target.name]: e.target.value
    });
  }

  handleSubmit(e) {
    this.props.onSubmit({
      Link: this.state.Link
    });

    this.setState({
      Link: "",
      showModal: false
    });

    e.preventDefault();
  }

  handleModalHide() {
    if(this.state.showModal) {
      this.setState({showModal: false});
    }
  }

  handleClick() {
    if(!this.state.showModal) {
      this.setState({showModal: true});
    }
  }

  render() {
    var form = (
      // jshint ignore:start
      <Form onSubmit={this.handleSubmit} body={
        <div>
        <Input name="Link" onChange={this.handleChange} value={this.state.Link}
        placeholder="Link" />

        <Submit text="Submit" />
        </div>
      } />
      // jshint ignore:end
    );

    // jshint ignore:start
    return (
      <div>
      <a onClick={this.handleClick} href="#">
      Subscribe To Feed
      </a>

      <Modal show={this.state.showModal} title="Subscribe To Feed"
      body={form} onHide={this.handleModalHide}/>
      </div>
    );
    // jshint ignore:end
  }
}

Subscribe.propTypes = {
  onSubmit: React.PropTypes.func.isRequired
};

module.exports = Subscribe;
