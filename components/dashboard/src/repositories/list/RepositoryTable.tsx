/**
 * Copyright (c) 2023 Gitpod GmbH. All rights reserved.
 * Licensed under the GNU Affero General Public License (AGPL).
 * See License.AGPL.txt in the project root for license information.
 */

import { FC } from "react";
import { TextInput } from "../../components/forms/TextInputField";
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@podkit/tables/Table";
import { LoaderIcon } from "lucide-react";
import { RepositoryListItem } from "./RepoListItem";
import { LoadingButton } from "@podkit/buttons/LoadingButton";
import { Configuration } from "@gitpod/public-api/lib/gitpod/v1/configuration_pb";
import { TextMuted } from "@podkit/typography/TextMuted";
import { Subheading } from "@podkit/typography/Headings";
import { cn } from "@podkit/lib/cn";

type Props = {
    configurations: Configuration[];
    searchTerm: string;
    hasNextPage: boolean;
    hasMoreThanOnePage: boolean;
    isSearching: boolean;
    isFetchingNextPage: boolean;
    onSearchTermChange: (val: string) => void;
    onLoadNextPage: () => void;
};

export const RepositoryTable: FC<Props> = ({
    searchTerm,
    configurations,
    hasNextPage,
    hasMoreThanOnePage,
    isSearching,
    isFetchingNextPage,
    onSearchTermChange,
    onLoadNextPage,
}) => {
    return (
        <>
            {/* Search/Filter bar */}
            <div className="flex flex-row flex-wrap justify-between items-center">
                <div className="flex flex-row flex-wrap gap-2 items-center">
                    {/* TODO: Add search icon on left and decide on pulling Inputs into podkit */}
                    <TextInput
                        className="w-80"
                        value={searchTerm}
                        onChange={onSearchTermChange}
                        placeholder="Search imported repositories"
                    />
                    {/* TODO: Add prebuild status filter dropdown */}
                </div>
            </div>
            <div className="relative w-full overflow-auto mt-2">
                {configurations.length > 0 ? (
                    <Table>
                        {/* TODO: Add sorting controls */}
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-52">Name</TableHead>
                                <TableHead hideOnSmallScreen>Repository</TableHead>
                                <TableHead className="w-32" hideOnSmallScreen>
                                    Created
                                </TableHead>
                                <TableHead className="w-24" hideOnSmallScreen>
                                    Prebuilds
                                </TableHead>
                                {/* Action column, loading status in header */}
                                <TableHead className="w-24 text-right">
                                    {isSearching && (
                                        <div className="flex flex-right justify-end items-center">
                                            {/* TODO: Make a LoadingIcon component */}
                                            <LoaderIcon
                                                className="animate-spin text-gray-500 dark:text-gray-300"
                                                size={20}
                                            />
                                        </div>
                                    )}
                                </TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {configurations.map((configuration) => {
                                return <RepositoryListItem key={configuration.id} configuration={configuration} />;
                            })}
                        </TableBody>
                    </Table>
                ) : (
                    <div
                        className={cn(
                            "w-full flex justify-center rounded-xl bg-gray-100 dark:bg-gray-800 px-4 py-10 animate-fade-in-fast",
                        )}
                    >
                        <Subheading className="max-w-md">No results found. Try adjusting your search terms.</Subheading>
                    </div>
                )}

                <div className="mt-4 mb-8 flex flex-row justify-center">
                    {hasNextPage ? (
                        <LoadingButton variant="secondary" onClick={onLoadNextPage} loading={isFetchingNextPage}>
                            Load more
                        </LoadingButton>
                    ) : (
                        hasMoreThanOnePage && <TextMuted>All repositories are loaded</TextMuted>
                    )}
                </div>
            </div>
        </>
    );
};
