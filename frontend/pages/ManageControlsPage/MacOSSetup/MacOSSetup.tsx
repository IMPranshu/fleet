import React, { useContext } from "react";
import PATHS from "router/paths";
import { Params } from "react-router/lib/Router";

import { AppContext } from "context/app";

import SideNav from "pages/admin/components/SideNav";
import Button from "components/buttons/Button/Button";
import PremiumFeatureMessage from "components/PremiumFeatureMessage";

import MAC_OS_SETUP_NAV_ITEMS from "./MacOSSetupNavItems";

const baseClass = "macos-setup";

interface ISetupEmptyState {
  router: any;
}

const SetupEmptyState = ({ router }: ISetupEmptyState) => {
  const onClickEmptyConnect = () => {
    router.push(PATHS.CONTROLS_MAC_SETTINGS);
  };

  return (
    <div className={`${baseClass}__empty-state`}>
      <h2>Setup experience for macOS hosts</h2>
      <p>Connect Fleet to the Apple Business Manager to get started.</p>
      <Button variant="brand" onClick={onClickEmptyConnect}>
        Connect
      </Button>
    </div>
  );
};

interface IMacOSSetupProps {
  params: Params;
  location: { search: string };
  router: any;
  teamIdForApi: number;
}

const MacOSSetup = ({
  params,
  location: { search: queryString },
  router,
  teamIdForApi,
}: IMacOSSetupProps) => {
  const { section } = params;
  const { isPremiumTier, config } = useContext(AppContext);

  const DEFAULT_SETTINGS_SECTION = MAC_OS_SETUP_NAV_ITEMS[0];

  const currentFormSection =
    MAC_OS_SETUP_NAV_ITEMS.find((item) => item.urlSection === section) ??
    DEFAULT_SETTINGS_SECTION;

  const CurrentCard = currentFormSection.Card;

  // TODO: uncomment when API done
  if (!config?.mdm.apple_bm_enabled_and_configured) {
    return <SetupEmptyState router={router} />;
  }

  return (
    <div className={baseClass}>
      <p>
        Customize the setup experience for hosts that automatically enroll to
        this team.
      </p>
      {!isPremiumTier ? (
        <PremiumFeatureMessage />
      ) : (
        <SideNav
          className={`${baseClass}__side-nav`}
          navItems={MAC_OS_SETUP_NAV_ITEMS.map((navItem) => ({
            ...navItem,
            path: navItem.path.concat(queryString),
          }))}
          activeItem={currentFormSection.urlSection}
          CurrentCard={
            <CurrentCard key={teamIdForApi} currentTeamId={teamIdForApi} />
          }
        />
      )}
    </div>
  );
};

export default MacOSSetup;
