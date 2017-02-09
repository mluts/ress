import React from "react";

class Feed extends React.Component {
  constructor(props) {
    super(props);
    this.selectFeed = this.selectFeed.bind(this);
    this.showContextMenu = this.showContextMenu.bind(this);
  }

  className() {
    var klass = ['feed'];
    if(this.props.data.Selected) {
      klass.push('selected');
    }
    if(this.props.data.Error) {
      klass.push("error");
    }
    return klass.join(' ');
  }

  selectFeed() {
    this.props.onSelectFeed(this.props.data);
  }

  showContextMenu(e) {
    e.preventDefault();
    e.stopPropagation();
    console.log(e);
    this.props.onContextMenu(this.props.data);
  }

  countUnread() {
    return this.props.data.Items.reduce((prev, cur) => {
      return cur.Unread ? prev + 1 : prev;
    }, 0);
  }

  render() {
    return (
      // jshint ignore:start
      <div title={this.props.data.Error} className="container-fluid" onContextMenu={this.showContextMenu}
          onClick={this.selectFeed} className={this.className()}>
      <div className="row">
        <div className="col-md-10">{this.props.data.Title}</div>
        <div className="col-md-2 unread pull-right">{this.countUnread()}</div>
      </div>
      </div>
      // jshint ignore:end
    );
  }
}

Feed.propTypes = {
  feed: React.PropTypes.object
};

class Feeds extends React.Component {
  render() {
    var feeds = this.props.data.map((f, i) => {
      return (
        // jshint ignore:start
        <li key={i.toString()}>
        <Feed onContextMenu={this.props.onFeedContextMenu} onSelectFeed={this.props.onSelectFeed} data={f} />
        </li>
        // jshint ignore:end
      );
    });
    // jshint ignore:start
    return <ul className="feeds">{feeds}</ul>;
    // jshint ignore:end
  }
}

Feeds.propTypes = {
  data: React.PropTypes.array,
  onSelectFeed: React.PropTypes.func
};

Feeds.Feed = Feed;
module.exports = Feeds;
