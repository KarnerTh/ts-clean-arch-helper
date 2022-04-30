const exampleValue = `export enum EnumWithValuesEntity {
  Draft = "draft",
  Published = "published",
}

export enum EnumWithoutValuesEntity {
  Monthly,
  Yearly,
}

export interface ExampleInterfaceEntity {
  name: string;
  description: string;
  enum1: EnumWithValuesEntity;
}

@Entity("ExampleEntity")
export class ExampleEntity {
  @PrimaryGeneratedColumn()
  id!: number;
  
  @Column("date")
  exampleDate!: Date;
  
  @Column("enum", { enum: EnumWithValuesEntity })
  enum1!: EnumWithValuesEntity;

  @Column("enum", { enum: EnumWithoutValuesEntity })
  enum2!: EnumWithoutValuesEntity;

  @ManyToOne(() => ExampleInterfaceEntity)
  @JoinColumn({ name: "item_id", referencedColumnName: "id" })
  singleItem!: ExampleInterfaceEntity;

  @OneToMany(() => ExampleInterfaceEntity, (example) => revision.item)
  multipleItems!: ExampleInterfaceEntity[];

  @Column("datetime")
  createdAt!: Date;

  @Column("datetime")
  updatedAt!: Date;
}`;
