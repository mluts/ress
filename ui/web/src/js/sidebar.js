var React = require('react'),
    _ = require('lodash');

class Sidebar extends React.Component {
  constructor(props) {
    super(props);
    _.bindAll(this, ['feedWasSelected']);
    this.state = { selectedFeedID: null };
  }

  feedWasSelected(feed) {
    this.setState({selectedFeedID: feed.ID});
    this.props.onFeedSelected(feed);
  }

  render() {
    var feeds = _.map(this.props.feeds, (f, i) => {
      // jshint ignore:start
      return(
        <li key={f.ID.toString()}>
          <Feed selected={f.ID == this.state.selectedFeedID}
                onSelected={this.feedWasSelected}
                data={f} />
        </li>
      );
      // jshint ignore:end
    });

    // jshint ignore:start
    return <ul>{feeds}</ul>;
    // jshint ignore:end
  }
}

Sidebar.propTypes = {
  feeds: React.PropTypes.array.isRequired,
  onFeedSelected: React.PropTypes.func.isRequired
};

class Feed extends React.Component {
  constructor(props) {
    super(props);
    _.bindAll(this, ['becomeSelected']);
  }

  becomeSelected() {
    this.props.onSelected(this.props.data);
  }

  render() {
    // jshint ignore:start
    return (
      <h3 className={this.props.selected ? "selected" : ""}
          onClick={this.becomeSelected}>{this.props.data.Title}</h3>
    );
    // jshint ignore:end
  }
}

Feed.propTypes = {
  data: React.PropTypes.object.isRequired,
  onSelected: React.PropTypes.func.isRequired
};

module.exports = Sidebar;
